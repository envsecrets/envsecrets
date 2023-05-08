package commons

type IntegrationType string

func (t *IntegrationType) IsValid() bool {
	for _, item := range AllowedIntegrations {
		if *t == item {
			return true
		}
	}
	return false
}

const (
	Github   IntegrationType = "github"
	Vercel   IntegrationType = "vercel"
	ASM      IntegrationType = "asm"
	CircleCI IntegrationType = "circle"
)

const (
	INTEGRATION_TYPE = "integration_type"
	INTEGRATION_ID   = "integration_id"
)

var (
	AllowedIntegrations = []IntegrationType{Github, Vercel, ASM, CircleCI}
)
