package gsm

import "github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"

type SetupOptions struct {
	Keys  map[string]interface{}
	OrgID string
}

type ListOptions struct {
	Credentials map[string]interface{}
	OrgID       string `json:"org_id"`
}

type SyncOptions struct {
	OrgID         string                 `json:"org_id"`
	Credentials   map[string]interface{} `json:"credentials"`
	EntityDetails map[string]interface{} `json:"entity_details"`
	Data          *keypayload.KPMap      `json:"data"`
}
