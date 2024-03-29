package auth

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	keyCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func SigninHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload SigninOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize a new HTTP client
	client := clients.NewNhostClient(&clients.NhostConfig{})

	//	Get the auth service.
	service := auth.GetService()

	//	Call the appropriate service handler.
	var response *auth.SigninResponse
	var err error
	if payload.Ticket == "" {
		response, err = service.SigninWithPassword(ctx, client, &auth.SigninWithPasswordOptions{
			Email:    payload.Email,
			Password: payload.Password,
		})
	} else {
		response, err = service.SigninWithMFA(ctx, client, &auth.SigninWithMFAOptions{
			Ticket: payload.Ticket,
			OTP:    payload.OTP,
		})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Login failed. Recheck your credentials.",
			Error:   err.Error(),
		})
	}

	if response.MFA != nil {
		return c.JSON(http.StatusOK, response)
	}

	if response.Session["accessToken"] == nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Login failed. Recheck your credentials.",
			Error:   "could not generate access token",
		})
	}

	//	Initialize a new GQL client with the user's access token.
	gqlClient := clients.NewGQLClient(&clients.GQLConfig{
		Authorization: "Bearer " + response.Session["accessToken"].(string),
		Type:          clients.HasuraClientType,
	})

	//	Extract and decrypt keys from user's session.
	pair, err := service.DecryptKeysFromSession(ctx, gqlClient, &auth.DecryptKeysFromSessionOptions{
		Session:  response.Session,
		Password: payload.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Login failed. Could not decrypt your keys.",
			Error:   err.Error(),
		})
	}

	//	Encrypt the sync key using the user's public key.
	var publicKey [32]byte
	copy(publicKey[:], pair.PublicKey)
	encryptedSyncKey, err := keys.SealAsymmetricallyAnonymous(pair.SyncKey, publicKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Login failed. Could not encrypt your sync key.",
			Error:   err.Error(),
		})
	}

	//	Include the decrypted keys in response.
	response.Keys = map[string]string{
		"publicKey":  base64.StdEncoding.EncodeToString(pair.PublicKey),
		"privateKey": base64.StdEncoding.EncodeToString(pair.PrivateKey),
		"syncKey":    base64.StdEncoding.EncodeToString(encryptedSyncKey),
	}

	return c.JSON(http.StatusOK, response)

	//return writeCookie(c, string(body))
}

func SignupHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload SignupOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Call the service handler.
	if err := auth.Signup(ctx, client, &auth.SignupOptions{
		Email:    payload.Email,
		Password: payload.Password,
		Name:     payload.Name,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to register the user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "Verification email sent to your inbox!",
		Data:    payload.Email,
	})
}

func UpdatePasswordHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload UpdatePasswordOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Extract the user's email from JWT
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*clients.Claims)

	//	Check whether user has keys.
	_, err := keys.GetByUserID(ctx, client, claims.Hasura.UserID)
	if err != nil {
		apiError := clients.ParseExternal(err)

		//	If key pair does not exist for this user,
		//	issue them a new key pair.
		if apiError.IsType(clients.ErrorTypeDoesNotExist) {

			//	Generate Key pair
			pair, err := keys.GenerateKeyPair(payload.NewPassword)
			if err != nil {
				return c.JSON(http.StatusBadRequest, &clients.APIResponse{
					Message: "Failed to issue a fresh key pair",
					Error:   err.Error(),
				})
			}

			//	Upload the keys to their cloud account.
			if err := keys.Create(ctx, client, &keyCommons.CreateOptions{
				PublicKey:    base64.StdEncoding.EncodeToString(pair.PublicKey),
				PrivateKey:   base64.StdEncoding.EncodeToString(pair.PrivateKey),
				ProtectedKey: base64.StdEncoding.EncodeToString(pair.ProtectedKey),
				Salt:         base64.StdEncoding.EncodeToString(pair.Salt),
			}); err != nil {
				return c.JSON(http.StatusBadRequest, &clients.APIResponse{
					Message: "Failed to issue a fresh key pair",
					Error:   err.Error(),
				})
			}
		} else {
			return c.JSON(http.StatusBadRequest, &clients.APIResponse{
				Error: err.Error(),
			})
		}
	}

	//	TODO: If a key-pair already exists,
	//	update the user's protection key.

	//	Initialize HTTP client
	httpClient := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service handler.
	if err := auth.UpdatePassword(ctx, httpClient, &auth.UpdatePasswordOptions{
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to update the password",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "password successfuly updated",
	})
}

func GenerateQRHandler(c echo.Context) error {

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize a new HTTP client
	client := clients.NewNhostClient(&clients.NhostConfig{
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the appropriate service handler.
	response, err := auth.GetService().GenerateTOTPQR(ctx, client)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to generate QR Code.",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func ToggleMFAHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload ToggleMFAOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize a new HTTP client
	client := clients.NewNhostClient(&clients.NhostConfig{
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Prepare service options.
	options := auth.ToggleMFAOptions{
		Code: payload.Code,
	}

	//	If it is a POST request, activate MFA.
	//	If it is a DELETE request, deactivate MFA.
	if c.Request().Method == http.MethodPost {
		options.ActiveMFAType = auth.TOTP
	}

	//	Call the appropriate service handler.
	err := auth.GetService().ToggleMFA(ctx, client, &options)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to toggle MFA.",
			Error:   err.Error(),
		})
	}

	message := "MFA Deactivated"
	if c.Request().Method == http.MethodPost {
		message = "MFA Activated"
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: message,
	})
}

func GetSyncKeyHandler(c echo.Context) error {

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize a new HTTP client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
		Type:          clients.HasuraClientType,
	})

	//	Call the appropriate service handler.
	syncKeyBytes, err := keys.GetSyncKey(ctx, client)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get sync key.",
			Error:   err.Error(),
		})
	}

	//	Decrypt the sync key using the server's own encryption key.
	syncKey, err := keys.OpenSymmetricallyByServer(syncKeyBytes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the sync key by the server.",
			Error:   err.Error(),
		})
	}

	//	Get the user's public key.
	publicKeyBytes, err := keys.GetPublicKey(ctx, client)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get user's public key.",
			Error:   err.Error(),
		})
	}

	//	Encrypt the sync key using the user's public key.
	var publicKey [32]byte
	copy(publicKey[:], publicKeyBytes)
	encryptedSyncKey, err := keys.SealAsymmetricallyAnonymous(syncKey, publicKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Login failed. Could not encrypt your sync key.",
			Error:   err.Error(),
		})
	}

	response := keyCommons.Key{
		SyncKey: base64.StdEncoding.EncodeToString(encryptedSyncKey),
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "Sync key fetched successfully.",
		Data:    &response,
	})
}

//	---	Helpers ---

func writeCookie(c echo.Context, value string) error {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = value
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return c.String(http.StatusOK, value)
}
