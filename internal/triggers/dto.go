package triggers

import "github.com/envsecrets/envsecrets/internal/auth"

type HasuraEventPayload struct {
	Event struct {
		Op               string            `json:"op"`
		SessionVariables auth.HasuraClaims `json:"session_variables"`
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
