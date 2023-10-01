package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/internal/auth/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	keyCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/nhost"
	"github.com/envsecrets/envsecrets/internal/organisations"
	organisationCommons "github.com/envsecrets/envsecrets/internal/organisations/commons"
	"github.com/envsecrets/envsecrets/internal/users"
	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
	"github.com/labstack/echo/v4"
)

func Signin(ctx context.ServiceContext, client *clients.HTTPClient, options *commons.SigninOptions) (*commons.SigninResponse, error) {

	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	//	Initialize a new request
	req, err := http.NewRequest(http.MethodPost, os.Getenv("NHOST_AUTH_URL")+"/v1/signin/email-password", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	//	Send the request to Nhost signin endpoint.
	var response commons.NhostSigninResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	//	Initialize a new GQL client with the user's access token.
	gqlClient := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   echo.HeaderAuthorization,
				Value: "Bearer " + response.Session["accessToken"].(string),
			},
		},
	})

	//	Extract the user's ID from the session.
	temp, err := json.Marshal(response.Session["user"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	var user userCommons.User
	if err := json.Unmarshal([]byte(temp), &user); err != nil {
		return nil, err
	}

	//	Fetch the keys of the user.
	ks, err := keys.GetByUserID(ctx, gqlClient, user.ID)
	if err != nil {
		return nil, err
	}

	//	Decode the keys.
	pair, err := ks.Decode()
	if err != nil {
		return nil, err
	}

	//	Decrypt the keys with user's password.
	if err := keys.DecryptPayload(pair, options.Password); err != nil {
		return nil, err
	}

	return &commons.SigninResponse{
		MFA:     response.MFA,
		Session: response.Session,
		Keys: map[string]string{
			"publicKey":  base64.StdEncoding.EncodeToString(pair.PublicKey),
			"privateKey": base64.StdEncoding.EncodeToString(pair.PrivateKey),
		},
	}, nil
}

func Signup(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SignupOptions) error {

	//	Signup on Nhost
	if err := nhost.Signup(ctx, &nhost.SignupOptions{
		Email:    options.Email,
		Password: options.Password,
		Options: map[string]interface{}{
			"displayName": options.Name,
		},
	}); err != nil {
		return errors.New(err.Message)
	}

	//	Fetch the user with their email.
	user, err := users.GetByEmail(ctx, client, options.Email)
	if err != nil {
		return err
	}

	//	Generate Key pair
	pair, err := keys.GenerateKeyPair(options.Password)
	if err != nil {
		return err
	}

	//	Upload the keys to their cloud account.
	if err := keys.CreateWithUserID(ctx, client, &keyCommons.CreateWithUserIDOptions{
		PublicKey:    base64.StdEncoding.EncodeToString(pair.PublicKey),
		PrivateKey:   base64.StdEncoding.EncodeToString(pair.PrivateKey),
		ProtectedKey: base64.StdEncoding.EncodeToString(pair.ProtectedKey),
		Salt:         base64.StdEncoding.EncodeToString(pair.Salt),
		UserID:       user.ID,
	}); err != nil {
		return err
	}

	//	Create a new `default` organisation for the new user.
	_, err = organisations.GetService().Create(ctx, client, &organisationCommons.CreateOptions{
		Name:   fmt.Sprintf("%s's Org", strings.Split(options.Name, " ")[0]),
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdatePassword(ctx context.ServiceContext, client *clients.HTTPClient, options *commons.UpdatePasswordOptions) error {

	body, err := options.Marshal()
	if err != nil {
		return err
	}

	//	Initialize a new request
	req, err := http.NewRequest(http.MethodPost, os.Getenv("NHOST_AUTH_URL")+"/user/password", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return client.Run(ctx, req, nil)
}
