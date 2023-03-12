package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/nhost"
	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
)

type LoginResponse struct {
	MFA struct {
		Ticket string `json:"ticket"`
	} `json:"mfa"`

	Session struct {
		AccessToken          string           `json:"accessToken"`
		AccessTokenExpiresIn int              `json:"accessTokenExpiresIn"`
		RefreshToken         string           `json:"refreshToken"`
		User                 userCommons.User `json:"user"`
	} `json:"session"`
}

func Login(payload map[string]interface{}) (*LoginResponse, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		configCommons.NHOST_AUTH_URL+"/signin/email-password",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(nhost.New(data).Message)
	}

	var response LoginResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

//	To logout the user, simply delete account config
func Logout() error {
	return config.GetService().Delete(configCommons.AccountConfig)
}

func IsLoggedIn() bool {
	return config.GetService().Exists(configCommons.AccountConfig)
}

func RefreshToken(payload map[string]interface{}) (*LoginResponse, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		configCommons.NHOST_AUTH_URL+"/token",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, nil
	}

	var response LoginResponse
	if err := json.Unmarshal(data, &response.Session); err != nil {
		return nil, err
	}

	return &response, nil
}

func RefreshAndSave() error {

	//	Fetch account configuration
	accountConfigPayload, err := config.GetService().Load(configCommons.AccountConfig)
	if err != nil {
		return err
	}

	accountConfig := accountConfigPayload.(*configCommons.Account)

	response, refreshErr := RefreshToken(map[string]interface{}{
		"refreshToken": accountConfig.RefreshToken,
	})

	if refreshErr != nil {
		return err
	}

	//	Save the refreshed account config
	refreshConfig := configCommons.Account{
		AccessToken:  response.Session.AccessToken,
		RefreshToken: response.Session.RefreshToken,
		User:         response.Session.User,
	}

	if err := config.GetService().Save(refreshConfig, configCommons.AccountConfig); err != nil {
		return err
	}

	return nil
}
