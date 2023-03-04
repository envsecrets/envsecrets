package clients

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
)

type HTTPClient struct {
	*http.Client
	BaseURL       string
	Authorization string
	Headers       []Header
	CustomHeaders []CustomHeader
}

type HTTPConfig struct {
	Type          ClientType
	BaseURL       string
	Authorization string
	Headers       []Header
	CustomHeaders []CustomHeader
}

func NewHTTPClient(config *HTTPConfig) *HTTPClient {

	var response HTTPClient
	response.Client = &http.Client{}

	if config == nil {
		return &response
	}

	response.Headers = config.Headers
	response.CustomHeaders = config.CustomHeaders
	response.BaseURL = config.BaseURL
	response.Authorization = config.Authorization

	switch config.Type {
	case GithubClientType:
		response.CustomHeaders = append(response.CustomHeaders, CustomHeader{
			Key:   string(AcceptHeader),
			Value: "application/vnd.github+json",
		})
	}

	return &response
}

func (c *HTTPClient) Run(ctx context.ServiceContext, req *http.Request) (*http.Response, *errors.Error) {

	/* 	//	Set headers
	   	for _, item := range c.Headers {
	   		switch item {
	   		case XHasuraAdminSecretHeader:
	   			req.Header.Add(string(item), os.Getenv(string(NHOST_ADMIN_SECRET)))
	   		}
	   	}
	*/

	//	Set Authorization Header
	if c.Authorization != "" {
		req.Header.Add(string(AuthorizationHeader), c.Authorization)
	}

	//	Set custom headers
	for _, item := range c.CustomHeaders {
		req.Header.Add(item.Key, item.Value)
	}

	//	Make the request
	response, err := c.Do(req)
	if err != nil {
		return nil, errors.New(err, "failed to send HTTP request", errors.ErrorTypeBadResponse, errors.ErrorSourceHTTP)
	}

	return response, nil
}
