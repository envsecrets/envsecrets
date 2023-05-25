package secrets

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/labstack/echo/v4"
)

func SetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.SetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Fetch the organisation using environment ID.
	organisation, err := organisations.GetByEnvironment(ctx, client, payload.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the organisation this environment is associated with",
			Error:   err.Error(),
		})
	}

	//	Get the server's copy of organisation's encryption key.
	var orgKey [32]byte
	orgKeyBytes, err := keys.GetOrgKeyServerCopy(ctx, organisation.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the organisation's encryption key'",
			Error:   err.Error(),
		})
	}
	copy(orgKey[:], orgKeyBytes)

	//	Since, the client is already sending base64 encoded values,
	//	mark the values as 'encoded'.
	vehicle := commons.Secret{
		Data: payload.Data,
	}
	vehicle.MarkEncoded()

	//	Encrypt the values with decrypted key
	if err := vehicle.Encrypt(orgKey); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to encrypt the secrets",
			Error:   err.Error(),
		})
	}

	//	Call the service function.
	secret, err := Set(ctx, client, &commons.SetSecretOptions{
		EnvID: payload.EnvID,
		Data:  vehicle.Data,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to set the secrets",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully set the secret",
		Data:    secret,
	})
}

func DeleteHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.DeleteRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service function.
	if _, err := Delete(ctx, client, &commons.DeleteSecretOptions{
		EnvID:   payload.EnvID,
		Key:     payload.Key,
		Version: payload.Version,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to delete the secret",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully delete the secret",
	})
}

func GetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.GetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
	})

	//	If the user has passed an authorization header,
	//	use that in GraphQL client.
	//	Else if they are authenticating using a token,
	//	it is safe to use the admin token.
	if c.Request().Header.Get(echo.HeaderAuthorization) != "" {
		client.Authorization = c.Request().Header.Get(echo.HeaderAuthorization)
	} else if c.Request().Header.Get(string(clients.TokenHeader)) != "" {
		client.Headers = append(client.Headers, clients.XHasuraAdminSecretHeader)
	} else {
		return echo.ErrUnauthorized
	}

	//	Override the env_id set by token middleware.
	if c.Get("env_id") != nil {
		payload.EnvID = c.Get("env_id").(string)
	}

	//	Call the service function.
	secret, err := Get(ctx, client, &commons.GetOptions{
		Key:     payload.Key,
		EnvID:   payload.EnvID,
		Version: payload.Version,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get the secrets",
			Error:   err.Error(),
		})
	}

	//	Fetch the organisation using environment ID.
	organisation, err := organisations.GetByEnvironment(ctx, client, payload.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the organisation this environment is associated with",
			Error:   err.Error(),
		})
	}

	//	Get the server's copy of organisation's encryption key.
	var orgKey [32]byte
	orgKeyBytes, err := keys.GetOrgKeyServerCopy(ctx, organisation.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the organisation's encryption key'",
			Error:   err.Error(),
		})
	}
	copy(orgKey[:], orgKeyBytes)

	//	Decrypt the values with decrypted key
	if err := secret.Decrypt(orgKey); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the secrets",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully got the secret",
		Data:    secret,
	})
}

func ListHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.ListRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{

			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
	})

	//	If the user has passed an authorization header,
	//	use that in GraphQL client.
	//	Else if they are authenticating using a token,
	//	it is safe to use the admin token.
	if c.Request().Header.Get(echo.HeaderAuthorization) != "" {
		client.Authorization = c.Request().Header.Get(echo.HeaderAuthorization)
	} else if c.Request().Header.Get(string(clients.TokenHeader)) != "" {
		client.Headers = append(client.Headers, clients.XHasuraAdminSecretHeader)
	} else {
		return echo.ErrUnauthorized
	}

	//	Override the env_id set by token middleware.
	if c.Get("env_id") != nil {
		payload.EnvID = c.Get("env_id").(string)
	}

	//	Call the service function.
	response, err := List(ctx, client, &payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the list of secrets",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully got the secret",
		Data:    response,
	})
}
