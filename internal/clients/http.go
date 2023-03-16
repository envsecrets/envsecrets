package clients

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/sirupsen/logrus"
)

type HTTPClient struct {
	*http.Client
	BaseURL       string
	Authorization string
	CustomHeaders []CustomHeader
	log           *logrus.Logger
}

type HTTPConfig struct {
	Type          ClientType
	BaseURL       string
	Authorization string
	Headers       []Header
	CustomHeaders []CustomHeader
	Logger        *logrus.Logger
}

func NewHTTPClient(config *HTTPConfig) *HTTPClient {

	var response HTTPClient
	response.Client = &http.Client{}

	if config == nil {
		return &response
	}

	response.CustomHeaders = config.CustomHeaders
	response.BaseURL = config.BaseURL
	response.Authorization = config.Authorization

	switch config.Type {
	case GithubClientType:
		response.CustomHeaders = append(response.CustomHeaders, CustomHeader{
			Key:   string(AcceptHeader),
			Value: "application/vnd.github+json",
		})
	case VaultClientType:
		response.CustomHeaders = append(response.CustomHeaders, CustomHeader{
			Key:   string(VaultTokenHeader),
			Value: os.Getenv("VAULT_ACCESS_TOKEN"),
		})
		response.CustomHeaders = append(response.CustomHeaders, CustomHeader{
			Key:   string(VaultNamespaceHeader),
			Value: "admin/default",
		})
	}

	if config.Logger != nil {
		response.log = config.Logger
	} else {
		response.log = logrus.New()
	}

	return &response
}

func (c *HTTPClient) Run(ctx context.ServiceContext, req *http.Request) (*http.Response, *errors.Error) {

	c.log.Debug("Sending HTTP request to: ", req.URL.String())

	//	Set content-type header
	req.Header.Set("content-type", "application/json")

	//	Set Authorization Header
	if c.Authorization != "" {
		req.Header.Set(string(AuthorizationHeader), c.Authorization)
	}

	//	Set custom headers
	for _, item := range c.CustomHeaders {
		req.Header.Set(item.Key, item.Value)
	}

	//	Backup the body in case it is required to re-run the request.
	body, err := req.GetBody()
	if err != nil {
		return nil, errors.New(err, "failed to create copy of request body", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	//	Make the request
	response, err := c.Do(req)
	if err != nil {
		return nil, errors.New(err, "failed to send HTTP request", errors.ErrorTypeBadResponse, errors.ErrorSourceHTTP)
	}

	//	If the request failed due to expired JWT,
	//	refresh the token and re-do the request.
	if response.StatusCode == 401 {

		c.log.Debug("Request failed due to expired token. Refreshing access token to try again.")

		//	Fetch account configuration
		accountConfigPayload, err := config.GetService().Load(configCommons.AccountConfig)
		if err != nil {
			return nil, errors.New(err, "failed to load account configuration", errors.ErrorTypeDoesNotExist, errors.ErrorSourceGo)
		}

		accountConfig := accountConfigPayload.(*configCommons.Account)

		response, refreshErr := auth.RefreshToken(map[string]interface{}{
			"refreshToken": accountConfig.RefreshToken,
		})

		if refreshErr != nil {
			return nil, errors.New(err, "failed to refresh auth token", errors.ErrorTypeBadResponse, errors.ErrorSourceNhost)
		}

		//	Save the refreshed account config
		refreshConfig := configCommons.Account{
			AccessToken:  response.Session.AccessToken,
			RefreshToken: response.Session.RefreshToken,
			User:         response.Session.User,
		}

		if err := config.GetService().Save(refreshConfig, configCommons.AccountConfig); err != nil {
			return nil, errors.New(err, "failed to save updated account configuration", errors.ErrorTypeInvalidAccountConfiguration, errors.ErrorSourceGo)
		}

		//	Update the authorization header in client.
		c.Authorization = "Bearer " + response.Session.AccessToken

		//	Re-set the body in the request, because it would have already been read once.
		req.Body = ioutil.NopCloser(body)

		return c.Run(ctx, req)
	}

	return response, nil
}
