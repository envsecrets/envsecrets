package commons

type JWTSecret struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}
