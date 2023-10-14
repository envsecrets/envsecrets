package clients

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"errors"

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

	//	Make the request
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	//	If the request failed due to expired JWT,
	//	refresh the token and re-do the request.
	if c.Type == HasuraClientType && resp.StatusCode == http.StatusForbidden {
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
