package netlify

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
)

func fetchAccounts(ctx context.ServiceContext, client *clients.HTTPClient) (*User, error) {

	req, err := http.NewRequest(http.MethodGet, "https://api.netlify.com/api/v1/user", nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := client.Run(ctx, req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
