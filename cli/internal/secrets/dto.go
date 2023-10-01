package secrets

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/dto"
)

type ListOptions struct {
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

type GetOptions struct {
	Key     string `json:"key"`
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

type DeleteOptions struct {
	Key     string `json:"key"`
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

type RemoteConfig struct {
	EnvironmentName string `json:"env_name"`
	ProjectID       string `json:"project_id"`
}

type Secret struct {
	dto.Secret
}

// Custom unmarsaller for the Secret struct.
func (s *Secret) UnmarshalJSON(data []byte) error {

	var secret dto.Secret
	if err := json.Unmarshal(data, &secret); err != nil {
		return err
	}

	*s = Secret{
		Secret: secret,
	}

	return nil
}

// Returns a boolean response subject to whether or not the secret has a non-empty environment ID.
func (s *Secret) IsRemote() bool {
	return s.EnvID != ""
}
