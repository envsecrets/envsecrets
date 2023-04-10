package actions

type HasuraActionPayload struct {
	Action struct {
		Name string `json:"name"`
	} `json:"action"`
	Input struct {
		Args interface{} `json:"args"`
	} `json:"input"`
}

type HasuraActionErrorResponse struct {
	Message    string        `json:"message"`
	Extensions []interface{} `json:"extensions"`
}
