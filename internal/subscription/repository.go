package subscription

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(context.Context, Subscription) error
	Get(context.Context, uuid.UUID) (Subscription, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context, ...Specification) ([]Subscription, error)
	Size(context.Context, ...Specification) (uint32, error)
}
