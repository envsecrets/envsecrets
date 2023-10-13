package payments

type Plan string

const (
	Monthly Plan = "monthly"
	Annual  Plan = "annual"
)

type CreateCheckoutSessionOptions struct {
	OrgID    string `query:"org_id"`
	Quantity int64  `query:"quantity"`
	Plan     Plan   `query:"plan"`
}
