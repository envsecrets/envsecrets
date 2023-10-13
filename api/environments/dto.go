package environments

import "github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"

type SyncWithPasswordOptions struct {
	EventIDs []string
	Password string `json:"password" validate:"required"`
	Version  *int   `json:"version,omitempty"`

	// Name of the secret to sync.
	Key string `json:"key,omitempty"`
}

type SyncOptions struct {
	EventIDs []string          `json:"event_ids,omitempty"`
	Pairs    *keypayload.KPMap `json:"pairs"`
}
