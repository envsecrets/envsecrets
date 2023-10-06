package integrations

type Type string

func (t *Type) IsValid() bool {
	for _, item := range AllowedIntegrations {
		if *t == item {
			return true
		}
	}
	return false
}

const (
	Github   Type = "github"
	Gitlab   Type = "gitlab"
	Vercel   Type = "vercel"
	ASM      Type = "asm"
	GSM      Type = "gsm"
	CircleCI Type = "circleci"
	Supabase Type = "supabase"
	Netlify  Type = "netlify"
)

var (
	AllowedIntegrations = []Type{Github, Gitlab, Vercel, ASM, CircleCI, GSM, Supabase, Netlify}
)
