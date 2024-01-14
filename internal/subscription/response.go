package subscription

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	DueAt     time.Time `json:"due_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubscriptionsResponse []SubscriptionResponse

type ListSubscriptionsResponse struct {
	Page          uint32                 `json:"page"`
	PageCount     uint32                 `json:"page_count"`
	PageSize      uint32                 `json:"page_size"`
	Size          uint32                 `json:"size"`
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
}

func SubscriptionResponseFromSubscription(subscription Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:        subscription.ID,
		Name:      subscription.Name,
		Type:      subscription.Type.String(),
		StartedAt: subscription.StartedAt,
		EndedAt:   subscription.EndedAt,
		DueAt:     subscription.DueAt,
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}
}

func SubscriptionsResponseFromSubscriptions(subscriptions Subscriptions) SubscriptionsResponse {
	var subscriptionsResponse SubscriptionsResponse

	for _, s := range subscriptions {
		subscriptionsResponse = append(subscriptionsResponse, SubscriptionResponseFromSubscription(s))
	}

	return subscriptionsResponse
}
