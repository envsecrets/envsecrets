package payments

import (
	"os"

	"github.com/stripe/stripe-go/v74"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_KEY")
}
