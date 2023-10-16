package subscriptions

import (
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

type Subscriptions []*Subscription

func (s Subscriptions) IsActiveAny() bool {
	for index := range s {
		if s[index].Status == "active" {
			return true
		}
	}
	return false
}
