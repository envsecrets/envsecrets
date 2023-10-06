package payments

type CreateCheckoutSessionOptions struct {
	OrgID    string `query:"org_id"`
	Quantity int64  `query:"quantity"`
}
