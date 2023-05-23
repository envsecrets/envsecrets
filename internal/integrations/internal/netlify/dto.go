package netlify

import (
	"encoding/json"

	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
)

type SetupOptions struct {
	Token string `json:"token"`
	OrgID string `json:"-"`
}

func (o *SetupOptions) toMap() (result map[string]interface{}) {
	payload, _ := json.Marshal(o)
	json.Unmarshal(payload, &result)
	return
}

type User struct {
	ID       string `json:"id"`
	UID      string `json:"uid"`
	FullName string `json:"full_name"`
}

type Site struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccountSlug string `json:"account_slug"`
}

type ListOptions struct {
	Credentials map[string]interface{}
}

type SyncOptions struct {
	Credentials   map[string]interface{} `json:"credentials"`
	EntityDetails map[string]interface{} `json:"entity_details"`
	Secret        secretCommons.Secret   `json:"secret"`
}
