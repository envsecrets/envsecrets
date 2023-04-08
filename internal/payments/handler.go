package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/payments/commons"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	"github.com/envsecrets/envsecrets/internal/users"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/webhook"
)

func CreateCheckoutSession(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.CreateCheckoutSessionOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Extract the user's email from JWT
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.Claims)

	user, err := users.Get(ctx, client, claims.Hasura.UserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &clients.APIResponse{
			Code:    http.StatusServiceUnavailable,
			Message: "failed to fetch user for this token",
			Error:   err.Error.Error(),
		})
	}

	quantity := payload.Quantity
	if quantity == 0 {
		quantity = 1
	}

	params := &stripe.CheckoutSessionParams{
		PhoneNumberCollection: &stripe.CheckoutSessionPhoneNumberCollectionParams{
			Enabled: stripe.Bool(true),
		},
		CustomerEmail:            stripe.String(user.Email),
		ClientReferenceID:        stripe.String(payload.OrgID),
		BillingAddressCollection: stripe.String("auto"),
		AllowPromotionCodes:      stripe.Bool(true),
		/* 		ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
		   			AllowedCountries: stripe.StringSlice([]string{
		   				"US",
		   				"IN",
		   			}),
		   		},
		*/
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				AdjustableQuantity: &stripe.CheckoutSessionLineItemAdjustableQuantityParams{
					Enabled: stripe.Bool(true),
					Minimum: stripe.Int64(1),
					Maximum: stripe.Int64(500),
				},

				// Provide the exact Price ID (for example, pr_1234) of the product you want to sell
				Price:    stripe.String(os.Getenv("STRIPE_PRICE_ID")),
				Quantity: stripe.Int64(quantity),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(os.Getenv("FE_URL") + "/access/members"),
		CancelURL:  stripe.String(os.Getenv("FE_URL") + "/access/members"),
	}

	stripe.Key = os.Getenv("STRIPE_KEY")
	s, er := session.New(params)
	if er != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to create a checkout session",
			Error:   er.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully created checkout sessions",
		Data: map[string]interface{}{
			"url": s.URL,
		},
	})
}

func CheckoutWebhook(c echo.Context) error {

	defer c.Request().Body.Close()

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error reading request body: %v\n", err))
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// You can find your endpoint's secret in your webhook settings
	event, err := webhook.ConstructEvent(body, c.Request().Header.Get("Stripe-Signature"), os.Getenv("STRIPE_CHECKOUT_WEBHOOK_ENDPOINT_SECRET"))

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error verifying webhook signature: %v\n", err))
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {

		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing webhook JSON: %v\n", err))
		}

		//	Insert a new subscription row.
		if _, err := subscriptions.Create(ctx, client, &subscriptions.CreateOptions{
			OrgID:          session.ClientReferenceID,
			SubscriptionID: session.Subscription.ID,
		}); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("%s: %s", err.Message, err.Error))
		}

	} else if event.Type == "customer.subscription.updated" {

		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing webhook JSON: %v\n", err))
		}

		//	Invite limit to update for the organisation
		quantity := event.Data.Object["quantity"].(float64)

		//	Check if this subscription already exists in our database,
		//	and it's value has merely been updated.
		if event.Data.PreviousAttributes["quantity"] != nil {
			previousQuantity := event.Data.PreviousAttributes["quantity"].(float64)

			//	Update the quantity
			quantity -= previousQuantity
		}

		//	Fetch the subscription using ID
		existingSubscription, er := subscriptions.GetBySubscriptionID(ctx, client, subscription.ID)
		if er != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Error fetching subscription from envsecrets database: %v\n", er))
		}

		//	Update the invite limit for the organisation.
		if err := organisations.UpdateInviteLimit(ctx, client, &organisations.UpdateInviteLimitOptions{
			ID:               existingSubscription.OrgID,
			IncrementLimitBy: int(quantity),
		}); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("%s: %s", err.Message, err.Error))
		}

		//	Update subscription status
		if _, err := subscriptions.Update(ctx, client, existingSubscription.ID, &subscriptions.UpdateOptions{
			Status: string(subscription.Status),
		}); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("%s: %s", err.Message, err.Error))
		}
	} else if event.Type == "customer.subscription.deleted" {

		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing webhook JSON: %v\n", err))
		}

		//	Delete the subscription from our database
		if err := subscriptions.DeleteBySubscriptionID(ctx, client, subscription.ID); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("%s: %s", err.Message, err.Error))
		}
	}

	return c.String(http.StatusOK, "webhook processed successfully")
}

/* func retrieveSessionItems(session_id string) (*[]stripe.LineItem, error) {

	req, err := http.NewRequest(http.MethodGet, "https://api.stripe.com/v1/checkout/sessions/"+session_id+"/line_items", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("STRIPE_KEY"), "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []stripe.LineItem `json:"data"`
	}
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
*/
