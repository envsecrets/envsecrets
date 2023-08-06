package gitlab

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
)

type EntityType string

const (
	ProjectType EntityType = "project"
	GroupType   EntityType = "group"
)

type SetupOptions struct {
	Token string
	OrgID string
	Code  string
}

type TokenRequestOptions struct {
	Code         string
	RedirectURI  string
	RefreshToken string
}

type TokenRefreshOptions struct {
	RefreshToken  string
	OrgID         string
	IntegrationID string
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type ListOptions struct {
	Credentials map[string]interface{}
	Type        EntityType
	Integration *commons.Integration
}

type SyncOptions struct {
	Credentials   map[string]interface{} `json:"credentials"`
	EntityDetails map[string]interface{} `json:"entity_details"`
	Data          *keypayload.KPMap      `json:"data"`
	IntegrationID string                 `json:"integration_id"`
	OrgID         string                 `json:"org_id"`
}

type ListProjectsResponse []Project

type Project struct {
	ID                int    `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Path              string `json:"path,omitempty"`
	PathWithNamespace string `json:"path_with_namespace,omitempty"`
	WebURL            string `json:"web_url,omitempty"`
}

type ListGroupsResponse []Group

type Group struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Path   string `json:"path,omitempty"`
	WebURL string `json:"web_url,omitempty"`
}

type CreateVariableOptions struct {
	ID       interface{}
	Variable Variable
}

type CreateVariableResponse struct {
	Message map[string]interface{} `json:"message"`
	Key     string                 `json:"key" form:"key"`
}

type Variable struct {
	Key              string `json:"key" form:"key"`
	Value            string `json:"value" form:"value"`
	Protected        bool   `json:"protected,omitempty" form:"protected,omitempty"`
	Masked           bool   `json:"masked,omitempty" form:"masked,omitempty"`
	EnvironmentScope string `json:"environment_scope,omitempty" form:"environment_scope,omitempty"`
}

func (v *Variable) Marshal() ([]byte, error) {
	return json.Marshal(v)
}
