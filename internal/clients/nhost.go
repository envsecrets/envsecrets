package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/sirupsen/logrus"
)

type NhostClient struct {
	*http.Client
	Authorization   string
	CustomHeaders   []CustomHeader
	log             *logrus.Logger
	ResponseHandler func(*http.Response) error
	BaseURL         string
}

type NhostConfig struct {
	Authorization   string
	Headers         []Header
	CustomHeaders   []CustomHeader
	Logger          *logrus.Logger
	ResponseHandler func(*http.Response) error
	BaseURL         string
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

	response.BaseURL = config.BaseURL
	if response.BaseURL == "" {
		response.BaseURL = os.Getenv("NHOST_AUTH_URL") + "/v1"
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

	//	Make the request
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	//	Unmarshal any errors in the response.
	var nhostError struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	if err := c.unmarshalResponse(resp, &nhostError); err != nil {
		return err
	}

	if nhostError.Error != "" {
		return fmt.Errorf("%v:%s:%s", nhostError.Status, nhostError.Error, nhostError.Message)
	}

	//	Unmarshal the remaining response.
	if err := c.unmarshalResponse(resp, &response); err != nil {
		return err
	}

	return nil
}

func (c *NhostClient) unmarshalResponse(resp *http.Response, structure interface{}) error {

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &structure); err != nil {
		return err
	}

	//	Rewrite the response body.
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}
