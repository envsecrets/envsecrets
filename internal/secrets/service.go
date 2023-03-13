package secrets

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/graphql"
)

func GenerateKey(ctx context.ServiceContext, path string, options commons.GenerateKeyOptions) *errors.Error {

	postBody, _ := options.Marshal()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/keys/"+path, bytes.NewBuffer(postBody))
	if err != nil {
		return errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	_, er := client.Run(ctx, req)
	if er != nil {
		return er
	}

	return nil
}

func DeleteKey(ctx context.ServiceContext, path string) *errors.Error {

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, os.Getenv("VAULT_ADDRESS")+"/v1/transit/keys/"+path, nil)
	if err != nil {
		return errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	_, er := client.Run(ctx, req)
	if er != nil {
		return er
	}

	return nil
}

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetSecretOptions) *errors.Error {

	//	If the secret type `ciphertext`,
	//	encrypt it from vault before saving the value.
	if options.Data.Payload.Type == commons.Ciphertext {

		postBody, _ := json.Marshal(options.GetVaultOptions())

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/encrypt/"+options.KeyPath, bytes.NewBuffer(postBody))
		if err != nil {
			return errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
		}

		client := clients.NewHTTPClient(&clients.HTTPConfig{
			Type: clients.VaultClientType,
		})

		resp, er := client.Run(ctx, req)
		if er != nil {
			return er
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(err, "failed to read response body", errors.ErrorTypeBadResponse, errors.ErrorSourceGo)
		}

		var response commons.VaultResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return errors.New(err, "failed to unmarshal set response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
		}

		//	Replace the secret value with ciphered version.
		options.Data.Payload.Value = response.Data.Ciphertext
	}

	//	Insert the encrypted secret in Hasura.
	if err := graphql.Set(ctx, client, options); err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DeleteSecretOptions) *errors.Error {

	//	Directly delete the key-value in Hasura.
	if err := graphql.Delete(ctx, client, options); err != nil {
		return err
	}

	return nil
}

func Get(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.GetResponse, *errors.Error) {

	//	Inittialize our secret data
	data := commons.Data{
		Key: options.Key,
	}

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {

		resp, err := graphql.GetByKeyByVersion(ctx, client, options)
		if err != nil {
			return nil, err
		}

		data.Payload = resp.Data[data.Key]

	} else {

		resp, err := graphql.GetByKey(ctx, client, options)
		if err != nil {
			return nil, err
		}

		data.Payload = resp.Data[data.Key]
		options.Version = &resp.Version
	}

	//	Only if the saved value was of type `ciphertext`,
	//	we have to descrypt the value.
	if data.Payload.Type == commons.Ciphertext {

		//	Decrypt the value from Vault.
		response, err := Decrypt(ctx, &commons.DecryptSecretOptions{
			Data:        data,
			KeyLocation: options.KeyPath,
			EnvID:       options.EnvID,
		})
		if err != nil {
			return nil, err
		}

		data.Payload.Value = response.Data.Plaintext
	}

	return &commons.GetResponse{
		Data: map[string]commons.Payload{
			data.Key: data.Payload,
		},
		Version: options.Version,
	}, nil
}

func GetAll(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.GetResponse, *errors.Error) {

	var data commons.GetResponse

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {

		resp, err := graphql.GetByVersion(ctx, client, options)
		if err != nil {
			return nil, err
		}

		data = *resp

	} else {

		resp, err := graphql.Get(ctx, client, options)
		if err != nil {
			return nil, err
		}

		data = *resp
	}

	//	Only if the saved value was of type `ciphertext`,
	//	we have to descrypt the value.
	for key, item := range data.Data {
		if item.Type == commons.Ciphertext {

			//	Decrypt the value from Vault.
			response, err := Decrypt(ctx, &commons.DecryptSecretOptions{
				Data: commons.Data{
					Key:     key,
					Payload: item,
				},
				KeyLocation: options.KeyPath,
				EnvID:       options.EnvID,
			})
			if err != nil {
				return nil, err
			}

			data.Data[key] = commons.Payload{
				Value: response.Data.Plaintext,
				Type:  item.Type,
			}
		}
	}

	return &data, nil
}

func Decrypt(ctx context.ServiceContext, options *commons.DecryptSecretOptions) (*commons.VaultResponse, *errors.Error) {

	postBody, _ := json.Marshal(options.GetVaultOptions())
	req, er := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/decrypt/"+options.KeyLocation, bytes.NewBuffer(postBody))
	if er != nil {
		return nil, errors.New(er, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	resp, err := client.Run(ctx, req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		return nil, errors.New(er, "failed to read response body", errors.ErrorTypeBadResponse, errors.ErrorSourceGo)
	}

	var response commons.VaultResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, errors.New(err, "failed to unmarshal set response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &response, nil
}
