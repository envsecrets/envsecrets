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
