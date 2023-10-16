package tokens

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	"github.com/envsecrets/envsecrets/internal/tokens"
)

// Abuse limit validation
//
// Before creating the token,
// we need to check whether the organisation's plan hasn't exceeded the abuse limit.
// If the number of list is greater than the allowed limit, proceed to check whether the organisation has an active subscription.
// Otherwise, approve the inputs and allow for creation of the project.
func validateInput(ctx context.ServiceContext, client *clients.GQLClient, env_id string) error {

	//	Fetch the organisation using environment ID.
	organisation, err := organisations.GetService().GetByEnvironment(ctx, client, env_id)
	if err != nil {
		return err
	}

	list, err := tokens.GetService().List(ctx, client, &tokens.ListOptions{
		EnvID: env_id,
	})
	if err != nil {
		return err
	}

	if len(list) >= FREE_TIER_LIMIT_NUMBER_OF_TOKENS {

		//	Validate whether the organisation an active premium subscription.
		//	We do this by fetching the subscriptions by the organisation ID.
		//	We then check if any subscription is active.
		subscriptions, err := subscriptions.GetService().GetByOrgID(ctx, client, organisation.ID)
		if err != nil {
			return err
		}

		//	If there are no subscriptions, or if even a single subscription is not active, return an error.
		if len(*subscriptions) == 0 {
			return clients.ErrBreachingAbuseLimit
		}

		active := subscriptions.IsActiveAny()
		if !active {
			return clients.ErrBreachingAbuseLimit
		}
	}

	return nil
}
