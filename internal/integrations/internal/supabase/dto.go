package supabase

import "github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"

type SetupOptions struct {
	Token string `json:"-"`
	OrgID string `json:"org_id"`
}

type Project struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Region         string `json:"region"`
}

type ListOptions struct {
	Credentials map[string]interface{}
}

type SyncOptions struct {
	Credentials   map[string]interface{} `json:"credentials"`
	EntityDetails map[string]interface{} `json:"entity_details"`
	Data          *keypayload.KPMap      `json:"data"`
}
