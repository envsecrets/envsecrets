package clients

import (
	"net/http"
	"os"

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
		BaseURL:       os.Getenv(string(NHOST_GRAPHQL_URL)),
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
