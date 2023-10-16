package clients

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/sirupsen/logrus"
)

type NhostClient struct {
	*clients.NhostClient
	log *logrus.Logger
}

type NhostConfig struct {
	Authorization string
	Logger        *logrus.Logger
	BaseURL       string
}

func NewNhostClient(config *NhostConfig) *NhostClient {

	client := clients.NewNhostClient(&clients.NhostConfig{
		BaseURL:       NHOST_AUTH_URL + "/v1",
		Authorization: config.Authorization,
	})

	response := NhostClient{
		NhostClient: client,
	}

	if config.Logger != nil {
		response.log = config.Logger
	} else {
		response.log = logrus.New()
	}

	return &response
}

func (c *NhostClient) Run(ctx context.ServiceContext, req *http.Request, response interface{}) error {

	c.log.Debug("[NhostClient] Request to: ", req.URL.String())

	return c.NhostClient.Run(ctx, req, &response)
}

func RefreshToken(payload map[string]interface{}) (*LoginResponse, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		NHOST_AUTH_URL+"/v1/token",
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
