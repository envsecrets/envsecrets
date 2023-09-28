package clients

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
		Data struct {
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
