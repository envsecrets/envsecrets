package tokens

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	keysCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	"github.com/envsecrets/envsecrets/internal/tokens"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func CreateHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload CreateOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with user's token
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Fetch the organisation using environment ID.
	organisation, err := organisations.GetService().GetByEnvironment(ctx, client, payload.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the organisation this environment is associated with",
			Error:   err.Error(),
		})
	}

	//	Initialize the tokens service.
	service := tokens.GetService()

	//
	//	Abuse limit validation
	//
	//	Before creating the token,
	//	we need to check whether the organisation's plan hasn't exceeded the abuse limit.
	//	If the number of list is greater than the allowed limit, proceed to check whether the organisation has an active subscription.
	//	Otherwise, approve the inputs and allow for creation of the project.
	list, err := service.List(ctx, client, &tokens.ListOptions{
		EnvID: payload.EnvID,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the tokens",
			Error:   err.Error(),
		})
	}

	if len(list) >= FREE_TIER_LIMIT_NUMBER_OF_TOKENS {

		//	Validate whether the organisation an active premium subscription.
		//	We do this by fetching the subscriptions by the organisation ID.
		//	We then check if any subscription is active.
		subscriptions, err := subscriptions.GetService().GetByOrgID(ctx, client, organisation.ID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &clients.APIResponse{
				Message: "failed to get the subscriptions",
				Error:   err.Error(),
			})
		}

		//	If there are no subscriptions, or if even a single subscription is not active, return an error.
		if len(*subscriptions) == 0 {
			return c.JSON(http.StatusBadRequest, &clients.APIResponse{
				Message: "Your current plan does not allow creating more tokens.",
				Error:   clients.ErrBreachingAbuseLimit.Error(),
			})
		}

		active := subscriptions.IsActiveAny()
		if !active {
			return c.JSON(http.StatusBadRequest, &clients.APIResponse{
				Message: "Your current plan does not allow creating more tokens.",
				Error:   clients.ErrBreachingAbuseLimit.Error(),
			})
		}
	}

	//	Extract the user's email from JWT
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.Claims)

	//	Decrypt and get the bytes of user's own copy of organisation's encryption key.
	orgKey, err := keys.DecryptMemberKey(ctx, client, claims.Hasura.UserID, &keysCommons.DecryptOptions{
		OrgID:    organisation.ID,
		Password: payload.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the organisation's encryption key. Maybe, entered password is invalid.",
			Error:   err.Error(),
		})
	}

	//	Create the token
	expiry, err := time.ParseDuration(payload.Expiry)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to parse expiry duration",
			Error:   err.Error(),
		})
	}

	token, err := service.Create(ctx, client, &tokens.CreateOptions{
		OrgKey: orgKey,
		EnvID:  payload.EnvID,
		Expiry: expiry,
		Name:   payload.Name,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to create the token",
			Error:   err.Error(),
		})
	}

	//	Encode the token
	hash := hex.EncodeToString(token)

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully generated token",
		Data: map[string]interface{}{
			"token": hash,
		},
	})
}
