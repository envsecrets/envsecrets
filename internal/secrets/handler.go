package secrets

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/labstack/echo/v4"
)

func SetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.SetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Fetch the organisation using environment ID.
	organisation, err := organisations.GetByEnvironment(ctx, client, payload.EnvID)
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to fetch the organisation this environment is associated with"),
			Error:   err.Message,
		})
	}

	//	Call the service function.
	secret, err := Set(ctx, client, &commons.SetSecretOptions{
		KeyPath:    organisation.ID,
		EnvID:      payload.EnvID,
		Data:       payload.Data,
		KeyVersion: payload.KeyVersion,
	})
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to set the secret"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully set the secret",
		Data:    secret,
	})
}

func MergeHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.MergeRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service function.
	secret, err := Merge(ctx, client, &commons.MergeSecretOptions{
		SourceEnvID:   payload.SourceEnvID,
		TargetEnvID:   payload.TargetEnvID,
		SourceVersion: payload.SourceVersion,
	})
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to merge the secrets"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully merged the secrets",
		Data:    secret,
	})
}

func DeleteHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.DeleteRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service function.
	if err := Delete(ctx, client, &commons.DeleteSecretOptions{
		EnvID: payload.EnvID,
		Key:   payload.Key,
	}); err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to delete the secret"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully delete the secret",
	})
}

func GetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.GetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

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

	//	Fetch the organisation using environment ID.
	organisation, err := organisations.GetByEnvironment(ctx, client, payload.EnvID)
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to fetch the organisation this environment is associated with"),
			Error:   err.Message,
		})
	}

	var response *commons.GetResponse

	//	If there is a specific key,
	//	pull the value only for that key.
	if payload.Key != "" {

		//	Call the service function.
		response, err = Get(ctx, client, &commons.GetSecretOptions{
			Key:     payload.Key,
			KeyPath: organisation.ID,
			EnvID:   payload.EnvID,
			Version: payload.Version,
		})
		if err != nil {
			return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
				Code:    err.Type.GetStatusCode(),
				Message: err.GenerateMessage("Failed to get the secret"),
				Error:   err.Message,
			})
		}

	} else {

		//	Else, pull all values.
		//	Call the service function.
		response, err = GetAll(ctx, client, &commons.GetSecretOptions{
			KeyPath: organisation.ID,
			EnvID:   payload.EnvID,
			Version: payload.Version,
		})
		if err != nil {
			return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
				Code:    err.Type.GetStatusCode(),
				Message: err.GenerateMessage("Failed to get the secret"),
				Error:   err.Message,
			})
		}
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully got the secret",
		Data:    response,
	})
}

func ListHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.ListRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

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
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to list the secret"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully got the secret",
		Data:    response,
	})
}
