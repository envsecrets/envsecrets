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

	req.Header.Set(string(commons.VAULT_TOKEN), os.Getenv(commons.VAULT_ROOT_TOKEN))

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return errors.New(err, "HTTP request failed to vault", errors.ErrorTypeRequestFailed, errors.ErrorSourceVault)
	}

	return nil
}

func DeleteKey(ctx context.ServiceContext, path string) *errors.Error {

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, os.Getenv("VAULT_ADDRESS")+"/v1/transit/keys/"+path, nil)
	if err != nil {
		return errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	req.Header.Set(string(commons.VAULT_TOKEN), os.Getenv(commons.VAULT_ROOT_TOKEN))

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return errors.New(err, "HTTP request failed to vault", errors.ErrorTypeRequestFailed, errors.ErrorSourceVault)
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

		req.Header.Set(string(commons.VAULT_TOKEN), os.Getenv(commons.VAULT_ROOT_TOKEN))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.New(err, "HTTP request failed to vault", errors.ErrorTypeRequestFailed, errors.ErrorSourceVault)
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

func Get(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetRequestOptions) (*commons.Secret, *errors.Error) {

	//	Inittialize our secret data
	data := commons.Data{
		Key: options.Key,
	}

	//	Initialize request options
	var getOptions = &commons.GetSecretOptions{
		EnvID: options.Path.Environment,
		Data:  data,
	}

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {

		getOptions.Version = options.Version

		resp, err := graphql.GetByKeyByVersion(ctx, client, getOptions)
		if err != nil {
			return nil, err
		}

		data.Payload = resp.Data[data.Key]

	} else {

		resp, err := graphql.GetByKey(ctx, client, getOptions)
		if err != nil {
			return nil, err
		}

		data.Payload = resp.Data[data.Key]
	}

	//	Save the returned encrypted value in our `get options`.
	getOptions.Data.Payload.Value = data.Payload.Value

	//	Only if the saved value was of type `ciphertext`,
	//	we have to descrypt the value.
	if data.Payload.Type == commons.Ciphertext {

		//	Decrypt the value from Vault.
		postBody, _ := json.Marshal(getOptions.GetVaultOptions())
		req, er := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/decrypt/"+options.Path.Organisation, bytes.NewBuffer(postBody))
		if er != nil {
			return nil, errors.New(er, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
		}

		req.Header.Set(string(commons.VAULT_TOKEN), os.Getenv(commons.VAULT_ROOT_TOKEN))

		resp, er := http.DefaultClient.Do(req)
		if er != nil {
			return nil, errors.New(er, "HTTP request failed to vault", errors.ErrorTypeRequestFailed, errors.ErrorSourceVault)
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

		data.Payload.Value = response.Data.Plaintext
		data.Payload.Type = commons.Plaintext
	}

	return &commons.Secret{
		Data: map[string]commons.Payload{
			options.Key: data.Payload,
		},
	}, nil
}
