package clients

import (
	"encoding/json"
	"io"
	"net/http"

	"errors"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/sirupsen/logrus"
)

type HTTPClient struct {
	*http.Client
	BaseURL         string
	Authorization   string
	CustomHeaders   []CustomHeader
	log             *logrus.Logger
	ResponseHandler func(*http.Response) error
	Type            ClientType
}

type HTTPConfig struct {
	Type            ClientType
	BaseURL         string
	Authorization   string
	Headers         []Header
	CustomHeaders   []CustomHeader
	Logger          *logrus.Logger
	ResponseHandler func(*http.Response) error
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
	response.ResponseHandler = config.ResponseHandler

	if config.Logger != nil {
		response.log = config.Logger
	} else {
		response.log = logrus.New()
	}

	return &response
}

func (c *HTTPClient) Run(ctx context.ServiceContext, req *http.Request, response interface{}) error {

	c.log.Debug("Sending request to: ", req.URL.String())

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

	req.Header.Get("content-type")

	//	Backup the body in case it is required to re-run the request.
	var body io.ReadCloser
	var err error
	if req.Body != nil {
		body, err = req.GetBody()
		if err != nil {
			return err
		}
	}

	//	Make the request
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	//	If the request failed due to expired JWT,
	//	refresh the token and re-do the request.
	if c.Type == HasuraClientType && resp.StatusCode == http.StatusUnauthorized {

		defer resp.Body.Close()

		result, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var errResponse APIResponse
		if err := json.Unmarshal(result, &errResponse); err != nil {
			return err
		}

		//	Only attempt the request again if the client has a JWT Authorization header attached.
		if c.Authorization == "" {
			return errors.New(errResponse.Message)
		}

		c.log.Debug("Refreshing access token to try again")

		//	Fetch account configuration
		accountConfigPayload, err := config.GetService().Load(configCommons.AccountConfig)
		if err != nil {
			return err
		}

		accountConfig := accountConfigPayload.(*configCommons.Account)

		authResponse, refreshErr := auth.RefreshToken(map[string]interface{}{
			"refreshToken": accountConfig.RefreshToken,
		})

		if refreshErr != nil {
			return err
		}

		//	Save the refreshed account config
		refreshConfig := configCommons.Account{
			AccessToken:  authResponse.Session.AccessToken,
			RefreshToken: authResponse.Session.RefreshToken,
			User:         authResponse.Session.User,
		}

		if err := config.GetService().Save(refreshConfig, configCommons.AccountConfig); err != nil {
			return err
		}

		//	Update the authorization header in client.
		c.Authorization = "Bearer " + authResponse.Session.AccessToken

		//	Re-set the body in the request, because it would have already been read once.
		if body != nil {
			req.Body = io.NopCloser(body)
		}

		return c.Run(ctx, req, response)

	} else if c.Type == HasuraClientType && resp.StatusCode == http.StatusForbidden {
		return errors.New("you do not have permission to perform this action")
	}

	if response != nil {

		defer resp.Body.Close()

		result, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(result, &response); err != nil {
			return err
		}
	}

	return nil
}
