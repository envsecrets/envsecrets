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
