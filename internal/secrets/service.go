package secrets

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/graphql"
)

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetSecretOptions) (*commons.Secret, *errors.Error) {

	for key, payload := range options.Data {

		//	If the secret type `ciphertext`,
		//	encrypt it from vault before saving the value.
		if payload.Type == commons.Ciphertext {

			postBody, _ := json.Marshal(map[string]interface{}{
				"plaintext":   payload.Value,
				"key_version": options.KeyVersion,
			})

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/encrypt/"+options.KeyPath, bytes.NewBuffer(postBody))
			if err != nil {
				return nil, errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
			}

			client := clients.NewHTTPClient(&clients.HTTPConfig{
				Type: clients.VaultClientType,
			})

			var response commons.VaultResponse
			if err := client.Run(ctx, req, &response); err != nil {
				return nil, err
			}

			//	Replace the secret value with ciphered version.
			payload.Value = response.Data.Ciphertext

			//	Update the map
			options.Data[key] = payload
		}
	}

	//	Insert the encrypted secret in Hasura.
	return graphql.Set(ctx, client, options)
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DeleteSecretOptions) *errors.Error {

	//	Directly delete the key-value in Hasura.
	if err := graphql.Delete(ctx, client, options); err != nil {
		return err
	}

	return nil
}

//	Cleanup entries from `secrets` table.
func Cleanup(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CleanupSecretOptions) *errors.Error {
	return graphql.Cleanup(ctx, client, options)
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
			Value:       data.Payload.Value,
			KeyLocation: options.KeyPath,
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
				Value:       item.Value,
				KeyLocation: options.KeyPath,
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

//	Fetches only the keys of a secret row.
func List(ctx context.ServiceContext, client *clients.GQLClient, options *commons.ListRequestOptions) (*commons.GetResponse, *errors.Error) {

	var data commons.GetResponse

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {

		resp, err := graphql.GetByVersion(ctx, client, &commons.GetSecretOptions{
			EnvID:   options.EnvID,
			Version: options.Version,
		})
		if err != nil {
			return nil, err
		}

		data = *resp

	} else {

		resp, err := graphql.Get(ctx, client, &commons.GetSecretOptions{
			EnvID: options.EnvID,
		})
		if err != nil {
			return nil, err
		}

		data = *resp
	}

	//	Remove the values from payload.
	//	Only keep the type.
	for key, item := range data.Data {
		data.Data[key] = commons.Payload{
			Type: item.Type,
		}
	}
	return &data, nil
}

//	Pulls all secret key-value pairs from the source environment,
//	and overwrites them in the target environment.
//	It creates a new secret version.
func Merge(ctx context.ServiceContext, client *clients.GQLClient, options *commons.MergeSecretOptions) (*commons.Secret, *errors.Error) {

	//	Fetch all key-value pairs of the source environment.
	sourceVariables, err := GetAll(ctx, client, &commons.GetSecretOptions{
		//	KeyPath: options.KeyPath,
		EnvID:   options.SourceEnvID,
		Version: options.SourceVersion,
	})
	if err != nil {
		return nil, err
	}

	//	Fetch all key-value pairs of the target environment.
	targetVariables, err := GetAll(ctx, client, &commons.GetSecretOptions{
		//	KeyPath: options.KeyPath,
		EnvID: options.TargetEnvID,
	})
	if err != nil {
		return nil, err
	}

	//	If the target variables is nil,
	//	then no pairs were fetched.
	if targetVariables.Data == nil {
		targetVariables.Data = make(map[string]commons.Payload)
	}

	//	Iterate through the target pairs,
	//	and overwrite the matching ones from the source pairs.
	for key, payload := range sourceVariables.Data {
		targetVariables.Data[key] = payload
	}

	//	Set the updated pairs in Hasura.
	return Set(ctx, client, &commons.SetSecretOptions{
		//	KeyPath: options.KeyPath,
		EnvID: options.TargetEnvID,
		Data:  targetVariables.Data,
	})
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

	var response commons.VaultResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
