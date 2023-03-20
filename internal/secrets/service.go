package secrets

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	internalErrors "errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/graphql"
)

//	This endpoint creates a new named encryption key of the specified type. The values set here cannot be changed after key creation.
//	Docs: https://developer.hashicorp.com/vault/api-docs/secret/transit#create-key
func GenerateKey(ctx context.ServiceContext, path string, options commons.GenerateKeyOptions) *errors.Error {

	postBody, _ := options.Marshal()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/keys/"+path, bytes.NewBuffer(postBody))
	if err != nil {
		return errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	var response commons.VaultResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return err
	}

	if len(response.Errors) != 0 {
		return errors.New(internalErrors.New(response.Errors[0].(string)), response.Errors[0].(string), errors.ErrorTypeBadResponse, errors.ErrorSourceVault)
	}

	return nil
}

//	This endpoint allows tuning configuration values for a given key. (These values are returned during a read operation on the named key.)
//	Docs: https://developer.hashicorp.com/vault/api-docs/secret/transit#update-key-configuration
func UpdateKeyConfiguration(ctx context.ServiceContext, path string, options commons.KeyConfigUpdateOptions) *errors.Error {

	postBody, _ := options.Marshal()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/keys/"+path+"/config", bytes.NewBuffer(postBody))
	if err != nil {
		return errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	return client.Run(ctx, req, nil)
}

//	This endpoint restores the backup as a named key. This will restore the key configurations and all the versions of the named key along with HMAC keys. The input to this endpoint should be the output of /backup endpoint.
//	Docs: https://developer.hashicorp.com/vault/api-docs/secret/transit#restore-key
func RestoreKey(ctx context.ServiceContext, path string, options commons.KeyRestoreOptions) *errors.Error {

	postBody, _ := options.Marshal()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/restore/"+path, bytes.NewBuffer(postBody))
	if err != nil {
		return errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	var response commons.VaultResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return err
	}

	if len(response.Errors) != 0 {
		return errors.New(internalErrors.New(response.Errors[0].(string)), response.Errors[0].(string), errors.ErrorTypeBadResponse, errors.ErrorSourceVault)
	}

	return nil
}

//	This endpoint returns a plaintext backup of a named key. The backup contains all the configuration data and keys of all the versions along with the HMAC key. The response from this endpoint can be used with the /restore endpoint to restore the key.
//	Docs: https://developer.hashicorp.com/vault/api-docs/secret/transit#backup-key
func BackupKey(ctx context.ServiceContext, path string) (*commons.KeyBackupResponse, *errors.Error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("VAULT_ADDRESS")+"/v1/transit/backup/"+path, nil)
	if err != nil {
		return nil, errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	var response commons.KeyBackupResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

//	This endpoint deletes a named encryption key. It will no longer be possible to decrypt any data encrypted with the named key. Because this is a potentially catastrophic operation, the deletion_allowed tunable must be set in the key's /config endpoint.
//	Docs: https://developer.hashicorp.com/vault/api-docs/secret/transit#delete-key
func DeleteKey(ctx context.ServiceContext, path string) *errors.Error {

	//	First, update key configuration. And make the key delete-able.
	if err := UpdateKeyConfiguration(ctx, path, commons.KeyConfigUpdateOptions{
		DeletionAllowed: true,
	}); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, os.Getenv("VAULT_ADDRESS")+"/v1/transit/keys/"+path, nil)
	if err != nil {
		return errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	return client.Run(ctx, req, nil)
}

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

	if err := graphql.Cleanup(ctx, client, options); err != nil {
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

//	Pulls all secret key-value pairs from the source environment,
//	and overwrites them in the target environment.
//	It creates a new secret version.
func Merge(ctx context.ServiceContext, client *clients.GQLClient, options *commons.MergeSecretOptions) (*commons.Secret, *errors.Error) {

	//	Fetch all key-value pairs of the source environment.
	sourceVariables, err := GetAll(ctx, client, &commons.GetSecretOptions{
		KeyPath: options.KeyPath,
		EnvID:   options.SourceEnvID,
		Version: options.SourceVersion,
	})
	if err != nil {
		return nil, err
	}

	//	Fetch all key-value pairs of the target environment.
	targetVariables, err := GetAll(ctx, client, &commons.GetSecretOptions{
		KeyPath: options.KeyPath,
		EnvID:   options.TargetEnvID,
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
		KeyPath: options.KeyPath,
		EnvID:   options.TargetEnvID,
		Data:    targetVariables.Data,
	})
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

	var response commons.VaultResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
