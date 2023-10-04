package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
)

type Service interface {
	ToggleMFA(context.ServiceContext, *clients.NhostClient, *commons.ToggleMFAOptions) error
	GenerateTOTPQR(context.ServiceContext, *clients.NhostClient) (*commons.GenerateQRResponse, error)
	SigninWithMFA(context.ServiceContext, *clients.NhostClient, *commons.SigninWithMFAOptions) (*commons.SigninResponse, error)
	SigninWithPassword(context.ServiceContext, *clients.NhostClient, *commons.SigninWithPasswordOptions) (*commons.SigninResponse, error)
	DecryptKeysFromSession(context.ServiceContext, *clients.GQLClient, *commons.DecryptKeysFromSessionOptions) (*keyCommons.Payload, error)
}

type DefaultService struct{}

// Remember: Passing a nil value to the "ActiveMFAType" option will deactivate MFA.
func (*DefaultService) ToggleMFA(ctx context.ServiceContext, client *clients.NhostClient, options *commons.ToggleMFAOptions) error {

	body, err := json.Marshal(options)
	if err != nil {
		return err
	}

	//	Initialize a new request
	req, err := http.NewRequest(http.MethodPost, client.BaseURL+"/v1/user/mfa", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	//	Send the request to Nhost signin endpoint.
	return client.Run(ctx, req, nil)
}

func (*DefaultService) GenerateTOTPQR(ctx context.ServiceContext, client *clients.NhostClient) (*commons.GenerateQRResponse, error) {

	//	Initialize a new request
	req, err := http.NewRequest(http.MethodGet, client.BaseURL+"/v1/mfa/totp/generate", nil)
	if err != nil {
		return nil, err
	}

	//	Send the request to Nhost signin endpoint.
	var response struct {
		Secret string `json:"totpSecret"`
		Image  string `json:"imageUrl"`
	}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &commons.GenerateQRResponse{
		Secret: response.Secret,
		Image:  response.Image,
	}, nil
}

func (*DefaultService) SigninWithMFA(ctx context.ServiceContext, client *clients.NhostClient, options *commons.SigninWithMFAOptions) (*commons.SigninResponse, error) {

	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	//	Initialize a new request
	req, err := http.NewRequest(http.MethodPost, client.BaseURL+"/v1/signin/mfa/totp", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	//	Send the request to Nhost signin endpoint.
	var response commons.NhostSigninResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &commons.SigninResponse{
		MFA:     response.MFA,
		Session: response.Session,
	}, nil
}

func (*DefaultService) SigninWithPassword(ctx context.ServiceContext, client *clients.NhostClient, options *commons.SigninWithPasswordOptions) (*commons.SigninResponse, error) {

	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	//	Initialize a new request
	req, err := http.NewRequest(http.MethodPost, client.BaseURL+"/v1/signin/email-password", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	//	Send the request to Nhost signin endpoint.
	var response commons.NhostSigninResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	//	Check whether the user has MFA enabled.
	if response.MFA != nil {
		return &commons.SigninResponse{
			MFA: response.MFA,
		}, nil
	}

	return &commons.SigninResponse{
		MFA:     response.MFA,
		Session: response.Session,
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
	req, err := http.NewRequest(http.MethodPost, client.BaseURL+"/user/password", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return client.Run(ctx, req, nil)
}

func (*DefaultService) DecryptKeysFromSession(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DecryptKeysFromSessionOptions) (*keyCommons.Payload, error) {

	//	Extract the user's ID from the session.
	temp, err := json.Marshal(options.Session["user"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	var user userCommons.User
	if err := json.Unmarshal([]byte(temp), &user); err != nil {
		return nil, err
	}

	//	Fetch the keys of the user.
	ks, err := keys.GetByUserID(ctx, client, user.ID)
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

	return pair, nil

}
