package subscription

import (
	"context"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Save(context.Context, Subscription) error
	Get(context.Context, uuid.UUID) (Subscription, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context, ...SubscriptionSpecification) ([]Subscription, error)
	Size(context.Context, ...SubscriptionSpecification) (uint32, error)
}
