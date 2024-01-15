package subscription

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MaybeTime time.Time

func (t MaybeTime) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte("null"), nil
	}

	return json.Marshal(tt)
}

type SubscriptionResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Fee       int32     `json:"fee"`
	Type      string    `json:"type"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   MaybeTime `json:"ended_at"`
	DueAt     time.Time `json:"due_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubscriptionsResponse []SubscriptionResponse

type ListSubscriptionsResponse struct {
	Page          uint32                `json:"page"`
	PageCount     uint32                `json:"page_count"`
	PageSize      uint32                `json:"page_size"`
	Size          uint32                `json:"size"`
	Subscriptions SubscriptionsResponse `json:"subscriptions"`
}

type CreateSubscriptionRequest struct {
	Name      string    `json:"name"`
	Fee       int32     `json:"fee"`
	Type      string    `json:"type"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	DueAt     time.Time `json:"due_at"`
}

type CreateSubscriptionResponse struct {
	Subscription SubscriptionResponse `json:"subscription"`
}

type GetSubscriptionResponse struct {
	Subscription SubscriptionResponse `json:"subscription"`
}

func NewSubscriptionResponse(subscription Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:        subscription.ID,
		Name:      subscription.Name,
		Fee:       subscription.Fee,
		Type:      subscription.Type.String(),
		StartedAt: subscription.StartedAt,
		EndedAt:   MaybeTime(subscription.EndedAt),
		DueAt:     subscription.DueAt,
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}
}

func NewSubscriptionsResponse(subscriptions Subscriptions) SubscriptionsResponse {
	subscriptionsResponse := SubscriptionsResponse{}

	for _, s := range subscriptions {
		subscriptionsResponse = append(subscriptionsResponse, NewSubscriptionResponse(s))
	}

	return subscriptionsResponse
}
