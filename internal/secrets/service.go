package secrets

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/client"
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

func Set(ctx context.ServiceContext, options *commons.SetOptions) *errors.Error {

	postBody, _ := json.Marshal(options.VaultOptions())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("VAULT_ADDRESS")+"/v1/transit/encrypt/"+options.Path.Organisation, bytes.NewBuffer(postBody))
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
	options.Secret.Value = response.Data.Ciphertext

	//	Insert the encrypted secret in Hasura.
	if err := graphql.Set(ctx, client.GRAPHQL_CLIENT, &commons.SetRequestOptions{
		EnvID:  options.Path.Environment,
		Secret: options.Secret,
	}); err != nil {
		return err
	}

	return nil
}

func Get(ctx context.ServiceContext, options *commons.GetRequestOptions) (*commons.Secret, *errors.Error) {

	//	Initialize new get request options.
	getOptions := commons.GetOptions{
		EnvID:  options.Path.Environment,
		Secret: commons.Secret{Key: options.Key},
	}

	//	Get the encrypted secret from Hasura.
	encryptedValue, err := graphql.GetByKey(ctx, client.GRAPHQL_CLIENT, &getOptions)
	if err != nil {
		return nil, err
	}

	//	Save the returned encrypted value in our `get options`.
	getOptions.Secret.Value = encryptedValue.Value

	//	Decrypt the value from Vault.
	postBody, _ := json.Marshal(getOptions.VaultOptions())
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

	//	Replace the secret value with plaintext version.
	var secret commons.Secret
	secret.Key = options.Key
	secret.Value = response.Data.Plaintext

	return &secret, nil
}

/*
func Get(ctx context.ServiceContext, key string, version *int) (*Secret, error) {

	//	Load the current project congi
	projectConfigData, er := config.GetService().Load(commons.ProjectConfig)
	if er != nil {
		panic(er.Error())
	}

	projectConfig := projectConfigData.(*commons.Project)

	//	Prepare body
	reqPayload := GetRequest{
		Key:     key,
		Version: version,
		Path: Path{
			Organisation: projectConfig.Organisation,
			Project:      projectConfig.Project,
			Environment:  projectConfig.Environment,
		},
	}

	body, err := reqPayload.Marshal()
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	addressPrefix := "/api/v1"
	req, err := http.NewRequest(
		http.MethodGet,
		os.Getenv("API")+addressPrefix+"/secrets/get",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	//	Set authorization header
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response APIResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	if response.Code != http.StatusOK {
		return nil, errors.New("failed to set the secret")
	}

	return &Secret{
		Key:   key,
		Value: response.Data,
	}, nil
}

func GetVersions(ctx context.ServiceContext, key string) (*Secret, error) {

	//	Load the current project congi
	projectConfigData, er := config.GetService().Load(commons.ProjectConfig)
	if er != nil {
		panic(er.Error())
	}

	projectConfig := projectConfigData.(*commons.Project)

	//	Prepare body
	reqPayload := GetRequest{
		Key: key,
		Path: Path{
			Organisation: projectConfig.Organisation,
			Project:      projectConfig.Project,
			Environment:  projectConfig.Environment,
		},
	}

	body, err := reqPayload.Marshal()
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	addressPrefix := "/api/v1"
	req, err := http.NewRequest(
		http.MethodGet,
		os.Getenv("API")+addressPrefix+"/secrets/get/versions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	//	Set authorization header
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response APIResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	if response.Code != http.StatusOK {
		return nil, errors.New("failed to set the secret")
	}

	fmt.Println(response.Data)

	return nil, nil
}
func List(ctx context.ServiceContext, version *int) (*[]Secret, error) {

	//	Load the current project congi
	projectConfigData, er := config.GetService().Load(commons.ProjectConfig)
	if er != nil {
		return nil, er
	}

	projectConfig := projectConfigData.(*commons.Project)

	//	Prepare body
	reqPayload := ListRequest{
		Version: version,
		Path: Path{
			Organisation: projectConfig.Organisation,
			Project:      projectConfig.Project,
			Environment:  projectConfig.Environment,
		},
	}

	body, err := reqPayload.Marshal()
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	addressPrefix := "/api/v1"
	req, err := http.NewRequest(
		http.MethodGet,
		os.Getenv("API")+addressPrefix+"/secrets/list",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	//	Set authorization header
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response APIResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	if response.Code != http.StatusOK {
		return nil, errors.New("failed to set the secret")
	}

	result, ok := response.Data.(map[string]interface{})
	if !ok {
		return nil, errors.New("failed type conversion for response data")
	}

	var secrets []Secret
	for key, value := range result {
		secrets = append(secrets, Secret{
			Key:   key,
			Value: value,
		})
	}

	return &secrets, nil
}
*/
