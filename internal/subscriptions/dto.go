package subscriptions

import (
	"encoding/json"
	"time"
)

type Subscription struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	OrgID          string    `json:"org_id,omitempty"`
	SubscriptionID string    `json:"subscription_id"`
	Status         Status    `json:"status"`
}

func (w *Subscription) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Subscription) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	OrgID          string `json:"org_id"`
	SubscriptionID string `json:"subscription_id"`
}

type UpdateOptions struct {
	Status string `json:"status"`
}

type ListOptions struct {
	OrgID  string `json:"org_id"`
	Status string `json:"status"`
}
