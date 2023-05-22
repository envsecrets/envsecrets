package netlify

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
)

func fetchAccounts(ctx context.ServiceContext, client *clients.HTTPClient) (*User, *errors.Error) {

	errMessage := "Failed to fetch accounts details from Netlify"

	req, er := http.NewRequest(http.MethodGet, "https://api.netlify.com/api/v1/user", nil)
	if er != nil {
		return nil, errors.New(er, errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	var user User
	if err := client.Run(ctx, req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
