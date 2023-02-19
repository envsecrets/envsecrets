package commons

type IntegrationType string

func (t *IntegrationType) IsValid() bool {
	return *t == Github ||
		*t == Vercel
}

const (
	Github IntegrationType = "github"
	Vercel IntegrationType = "vercel"
)

const (
	INTEGRATION_TYPE = "integration_type"
)
