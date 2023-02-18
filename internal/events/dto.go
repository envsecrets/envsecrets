package events

type HasuraEventPayload struct {
	Event struct {
		Op               string       `json:"op"`
		SessionVariables HasuraClaims `json:"session_variables"`
		Data             struct {
			Old interface{} `json:"old"`
			New interface{} `json:"new"`
		} `json:"data"`
	} `json:"event"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type HasuraClaims struct {
	AllowedRoles []string `json:"x-hasura-allowed-roles,omitempty"`
	DefaultRole  string   `json:"x-hasura-default-role,omitempty"`
	UserID       string   `json:"x-hasura-user-id,omitempty"`
	IsAnonymous  string   `json:"x-hasura-user-is-anonymous,omitempty"`
}
