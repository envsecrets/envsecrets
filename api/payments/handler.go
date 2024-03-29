package payments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	"github.com/envsecrets/envsecrets/internal/users"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/webhook"
)

func CreateCheckoutSessionHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload CreateCheckoutSessionOptions
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

	//	Extract the user's email from JWT
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.Claims)

	user, err := users.Get(ctx, client, claims.Hasura.UserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &clients.APIResponse{
			Message: "failed to fetch user for this token",
			Error:   err.Error(),
		})
	}

	//	Should be minimum 1 quantity.
	quantity := payload.Quantity
	if quantity == 0 {
		quantity = 1
	}

	//	Choose the monthly plan by default.
	if payload.Plan == "" {
		payload.Plan = Monthly
	}

	//	Initialize the appropriate Stripe Price ID.
	var priceID string
	switch payload.Plan {
	case Monthly:
		priceID = os.Getenv("STRIPE_MONTHLY_PLAN_PRICE_ID")
	case Annual:
		priceID = os.Getenv("STRIPE_ANNUAL_PLAN_PRICE_ID")
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
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(quantity),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(os.Getenv("FE_URL") + "/billing"),
		CancelURL:  stripe.String(os.Getenv("FE_URL") + "/billing"),
	}

	stripe.Key = os.Getenv("STRIPE_KEY")
	s, err := session.New(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to create a checkout session",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully created checkout sessions",
		Data: map[string]interface{}{
			"url": s.URL,
		},
	})
}

func CheckoutWebhookHandler(c echo.Context) error {

	defer c.Request().Body.Close()

	body, err := io.ReadAll(c.Request().Body)
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
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Get the subscription service.
	service := subscriptions.GetService()

	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {

		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing webhook JSON: %v\n", err))
		}

		//	Insert a new subscription row.
		if _, err := service.Create(ctx, client, &subscriptions.CreateOptions{
			OrgID:          session.ClientReferenceID,
			SubscriptionID: session.Subscription.ID,
		}); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("failed to register the new subscription: %s", err.Error()))
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
		existingSubscription, err := service.GetBySubscriptionID(ctx, client, subscription.ID)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Error fetching subscription from envsecrets database: %v\n", err))
		}

		//	Update the invite limit for the organisation.
		if err := organisations.GetService().UpdateInviteLimit(ctx, client, &organisations.UpdateInviteLimitOptions{
			ID:               existingSubscription.OrgID,
			IncrementLimitBy: int(quantity),
		}); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("failed to update org's invite limit: %s", err.Error()))
		}

		//	Update subscription status
		if _, err := service.Update(ctx, client, existingSubscription.ID, &subscriptions.UpdateOptions{
			Status: string(subscription.Status),
		}); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("%s: %s", "failed to update subscription status", err.Error()))
		}
	} else if event.Type == "customer.subscription.deleted" {

		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing webhook JSON: %v\n", err))
		}

		//	Delete the subscription from our database
		if err := service.DeleteBySubscriptionID(ctx, client, subscription.ID); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("%s: %s", "failed to delete the subscription record from db", err.Error()))
		}
	}

	return c.String(http.StatusOK, "webhook processed successfully")
}
