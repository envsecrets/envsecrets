package secrets

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/labstack/echo/v4"
)

func SetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.SetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
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
	secret, err := Set(ctx, client, &commons.SetSecretOptions{
		KeyPath:    payload.OrgID,
		EnvID:      payload.EnvID,
		Data:       payload.Data,
		KeyVersion: payload.KeyVersion,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to set the secret",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully set the secret",
		Data:    secret,
	})
}

func MergeHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.MergeRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
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
		KeyPath:     payload.OrgID,
		SourceEnvID: payload.SourceEnvID,
		TargetEnvID: payload.TargetEnvID,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to merge secrets",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully merged the secrets",
		Data:    secret,
	})
}

func DeleteHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.DeleteRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
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
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to set the secret",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully set the secret",
	})
}

func GetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.GetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
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

	var response *commons.GetResponse
	var err *errors.Error

	//	If there is a specific key,
	//	pull the value only for that key.
	if payload.Key != "" {

		//	Call the service function.
		response, err = Get(ctx, client, &commons.GetSecretOptions{
			Key:     payload.Key,
			KeyPath: payload.OrgID,
			EnvID:   payload.EnvID,
			Version: payload.Version,
		})
		if err != nil {
			return c.JSON(http.StatusBadRequest, &commons.APIResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to get the secret",
				Error:   err.Message,
			})
		}

	} else {

		//	Else, pull all values.
		//	Call the service function.
		response, err = GetAll(ctx, client, &commons.GetSecretOptions{
			KeyPath: payload.OrgID,
			EnvID:   payload.EnvID,
			Version: payload.Version,
		})
		if err != nil {
			return c.JSON(http.StatusBadRequest, &commons.APIResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to get the secret",
				Error:   err.Message,
			})
		}
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully got the secret",
		Data:    response,
	})
}

func KeyBackupHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.KeyBackupRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Call the service function.
	response, err := BackupKey(ctx, payload.OrgID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to set the secret",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully generated key plaintext backup",
		Data:    response.Data,
	})
}

func KeyRestoreHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.KeyRestoreRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Call the service function.
	err := RestoreKey(ctx, payload.OrgID, commons.KeyRestoreOptions{
		Backup: payload.Backup,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to set the secret",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully restored key from plaintext backup",
	})
}
