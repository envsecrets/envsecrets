package clients

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/sirupsen/logrus"
)

type NhostClient struct {
	*http.Client
	Authorization   string
	CustomHeaders   []CustomHeader
	log             *logrus.Logger
	ResponseHandler func(*http.Response) error
}

type NhostConfig struct {
	Authorization   string
	Headers         []Header
	CustomHeaders   []CustomHeader
	Logger          *logrus.Logger
	ResponseHandler func(*http.Response) error
}

func NewNhostClient(config *NhostConfig) *NhostClient {

	var response NhostClient
	response.Client = &http.Client{}
	response.log = logrus.New()

	if config == nil {
		return &response
	}

	if config.Logger != nil {
		response.log = config.Logger
	}

	response.CustomHeaders = config.CustomHeaders
	response.Authorization = config.Authorization
	response.ResponseHandler = config.ResponseHandler

	return &response
}

func (c *NhostClient) Run(ctx context.ServiceContext, req *http.Request, response interface{}) error {

	c.log.Debug("[NhostClient] Request to: ", req.URL.String())

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
