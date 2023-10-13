package nhost

import "github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"

type SetupOptions struct {
	Token string
	OrgID string `json:"org_id"`
}

type ListOptions struct {
	Credentials map[string]interface{}
}

type SyncOptions struct {
	OrgID         string                 `json:"org_id"`
	Credentials   map[string]interface{} `json:"credentials"`
	EntityDetails map[string]interface{} `json:"entity_details"`
	Data          *keypayload.KPMap      `json:"data"`
}
