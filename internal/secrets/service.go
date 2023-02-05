package secrets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/config"
	"github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/context"
)

func Set(ctx context.ServiceContext, data *Secret) error {

	//	Load the current project congi
	projectConfigData, er := config.GetService().Load(commons.ProjectConfig)
	if er != nil {
		panic(er.Error())
	}

	projectConfig := projectConfigData.(*commons.Project)

	//	Prepare body
	reqPayload := SetRequest{
		Secret: *data,
		Path: Path{
			Organisation: projectConfig.Organisation,
			Project:      projectConfig.Project,
			Environment:  projectConfig.Environment,
		},
	}

	body, err := reqPayload.Marshal()
	if err != nil {
		return err
	}

	client := http.Client{}

	addressPrefix := "/api/v1"
	req, err := http.NewRequest(
		http.MethodPost,
		os.Getenv("API")+addressPrefix+"/secrets/set",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	//	Set authorization header
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response APIResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return err
	}

	if response.Code != http.StatusOK {
		return errors.New("failed to set the secret")
	}

	return nil
}

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
