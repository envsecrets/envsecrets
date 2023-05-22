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
	Gitlab   IntegrationType = "gitlab"
	Vercel   IntegrationType = "vercel"
	ASM      IntegrationType = "asm"
	GSM      IntegrationType = "gsm"
	CircleCI IntegrationType = "circle"
	Supabase IntegrationType = "supabase"
	Netlify  IntegrationType = "netlify"
)

const (
	INTEGRATION_TYPE = "integration_type"
	INTEGRATION_ID   = "integration_id"
)

var (
	AllowedIntegrations = []IntegrationType{Github, Gitlab, Vercel, ASM, CircleCI, GSM, Supabase, Netlify}
)
