package keys

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	internalErrors "errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/keys/commons"
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

	return client.Run(ctx, req, nil)
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

//	This endpoint returns the named key. The keys object shows the value of the key for each version. If version is specified, the specific version will be returned. If latest is provided as the version, the current key will be provided. Depending on the type of key, different information may be returned. The key must be exportable to support this operation and the version must still be valid.
//	Docs: https://developer.hashicorp.com/vault/api-docs/secret/transit#export-key
func ExportKey(ctx context.ServiceContext, options *commons.KeyExportOptions) (*commons.KeyExportResponse, *errors.Error) {

	//	Export the latest key if no version is specified.
	if options.Version == "" {
		options.Version = "latest"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/transit/export/%s/%s/%s", os.Getenv("VAULT_ADDRESS"), options.Type, options.Name, options.Version), nil)
	if err != nil {
		return nil, errors.New(err, "failed to create HTTP request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VaultClientType,
	})

	var response commons.KeyExportResponse
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
