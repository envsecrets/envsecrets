package keys

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/labstack/echo/v4"
)

func KeyBackupHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.KeyBackupRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{

			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Call the service function.
	response, err := BackupKey(ctx, payload.OrgID)
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{

			Message: err.GenerateMessage("Failed to generate key backup"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{

		Message: "successfully generated key plaintext backup",
		Data:    response.Data,
	})
}

func KeyRestoreHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.KeyRestoreRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{

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
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{

			Message: err.GenerateMessage("Failed to restore the key"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{

		Message: "successfully restored key from plaintext backup",
	})
}
