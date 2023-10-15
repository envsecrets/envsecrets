package clients

import (
	"github.com/golang-jwt/jwt/v4"
)

type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Response struct {
	StatusCode int
	Data       interface{}
}

type HasuraTriggerPayload struct {
	Event struct {
		Op               string       `json:"op"`
		SessionVariables HasuraClaims `json:"session_variables"`
		Data             struct {
			New interface{} `json:"new"`
			Old interface{} `json:"old"`
		} `json:"data"`
	} `json:"event"`
}

type HasuraActionRequestPayload struct {
	Action struct {
		Name string `json:"name"`
	} `json:"action"`
	Input struct {
		Args interface{} `json:"args"`
	} `json:"input"`
}

type HasuraInputValidationPayload struct {
	Data struct {
		Input interface{} `json:"input"`
	} `json:"data"`
}

type HasuraClaims struct {
	AllowedRoles []string `json:"x-hasura-allowed-roles,omitempty"`
	DefaultRole  string   `json:"x-hasura-default-role,omitempty"`
	UserID       string   `json:"x-hasura-user-id,omitempty"`
	UserEmail    string   `json:"x-hasura-user-email,omitempty"`
	IsAnonymous  string   `json:"x-hasura-user-is-anonymous,omitempty"`
}

type Claims struct {
	Hasura HasuraClaims `json:"https://hasura.io/jwt/claims"`
	jwt.RegisteredClaims
}

type HasuraActionResponse struct {
	Message    string                          `json:"message"`
	Extensions HasuraActionsResponseExtensions `json:"extensions,omitempty"`
}

type HasuraActionsResponseExtensions struct {
	Code  int                    `json:"code"`
	Error error                  `json:"error,omitempty"`
	Data  map[string]interface{} `json:"data,omitempty"`
}
