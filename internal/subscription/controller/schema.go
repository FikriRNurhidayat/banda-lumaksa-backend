package subscription_controller

import (
	"encoding/json"
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/schema"
	"github.com/google/uuid"

	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/entity"
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
	common_schema.PaginationResponse
	Subscriptions SubscriptionsResponse `json:"subscriptions"`
}

type SubscriptionRequest struct {
	Name      string    `json:"name"`
	Fee       int32     `json:"fee"`
	Type      string    `json:"type"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	DueAt     time.Time `json:"due_at"`
}

type CreateSubscriptionRequest struct {
	Subscription SubscriptionRequest `json:"subscription"`
}

type CreateSubscriptionResponse struct {
	Subscription SubscriptionResponse `json:"subscription"`
}

type GetSubscriptionResponse struct {
	Subscription SubscriptionResponse `json:"subscription"`
}

func NewSubscriptionResponse(subscription subscription_entity.Subscription) SubscriptionResponse {
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

func NewSubscriptionsResponse(subscriptions subscription_entity.Subscriptions) SubscriptionsResponse {
	subscriptionsResponse := SubscriptionsResponse{}

	for _, s := range subscriptions {
		subscriptionsResponse = append(subscriptionsResponse, NewSubscriptionResponse(s))
	}

	return subscriptionsResponse
}
