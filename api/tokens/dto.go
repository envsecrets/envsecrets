package tokens

type CreateOptions struct {
	Password string `json:"password"`
	EnvID    string `json:"env_id"`
	Expiry   string `json:"expiry"`
	Name     string `json:"name,omitempty"`
}
