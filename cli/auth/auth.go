package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/internal/nhost"
	"github.com/envsecrets/envsecrets/internal/users"
)

type LoginResponse struct {
	MFA struct {
		Ticket string `json:"ticket"`
	} `json:"mfa"`

	Session NhostSession `json:"session"`
}

type NhostSession struct {
	AccessToken          string     `json:"accessToken"`
	AccessTokenExpiresIn int        `json:"accessTokenExpiresIn"`
	RefreshToken         string     `json:"refreshToken"`
	User                 users.User `json:"user"`
}

func Login(payload map[string]interface{}) (*LoginResponse, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		configCommons.NHOST_AUTH_URL+"/v1"+"/signin/email-password",
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

	data, err := io.ReadAll(res.Body)
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

// To logout the user, simply delete account config
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
		configCommons.NHOST_AUTH_URL+"/v1/token",
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

	data, err := io.ReadAll(res.Body)
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
