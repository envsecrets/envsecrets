package graphql

import "time"

type Integration struct {
	ID             string    `json:"id,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	UserID         string    `json:"user_id,omitempty"`
	OrgID          string    `json:"org_id"`
	InstallationID string    `json:"installation_id"`
	Type           string    `json:"type"`
	Credentials    string    `json:"credentials,omitempty"`
}

type UpdateDetailsOptions struct {
	ID            string                 `json:"id"`
	EntityDetails map[string]interface{} `json:"entity_details"`
}

type UpdateCredentialsOptions struct {
	ID          string `json:"id"`
	Credentials string `json:"credentials"`
}

type ListIntegrationFilters struct {
	OrgID string `json:"org_id"`
	Type  string `json:"type"`
}

type AddIntegrationOptions struct {

	//	Global
	OrgID          string `json:"org_id"`
	InstallationID string `json:"installation_id"`
	Type           string `json:"type"`

	//	Especially for Vercel
	Credentials string                 `json:"credentials,omitempty"`
	Scope       map[string]interface{} `json:"scope,omitempty"`
}
