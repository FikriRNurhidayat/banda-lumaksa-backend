package subscription

import (
	"context"
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/values"
	"github.com/google/uuid"
)

type CreateSubscriptionUseCase interface {
	Call(ctx context.Context, params *CreateSubscriptionParams) (*CreateSubscriptionResult, error)
}

type CreateSubscriptionParams struct {
	Name      string
	Fee       int32
	Type      Type
	StartedAt time.Time
	EndedAt   time.Time
	DueAt     time.Time
}

type CreateSubscriptionResult struct {
	Subscription Subscription
}

type CreateSubscriptionUseCaseImpl struct {
	subscriptionRepository Repository
}

func (u *CreateSubscriptionUseCaseImpl) Call(ctx context.Context, params *CreateSubscriptionParams) (*CreateSubscriptionResult, error) {
	s := Subscription{
		ID:        uuid.New(),
		Name:      params.Name,
		Fee:       params.Fee,
		Type:      params.Type,
		StartedAt: params.StartedAt,
		EndedAt:   params.EndedAt,
		DueAt:     params.DueAt,
	}

	if s.DueAt == values.EmptyTime {
		s.DueAt = u.ComputeDueAt(s)
	}

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	if err := u.subscriptionRepository.Save(ctx, s); err != nil {
		return nil, err
	}

	return &CreateSubscriptionResult{
		Subscription: s,
	}, nil
}

func NewCreateSubscriptionUseCase(subscriptionRepository Repository) CreateSubscriptionUseCase {
	return &CreateSubscriptionUseCaseImpl{
		subscriptionRepository: subscriptionRepository,
	}
}

func (u *CreateSubscriptionUseCaseImpl) ComputeDueAt(s Subscription) time.Time {
	switch s.Type {
	case Daily:
		return s.StartedAt.Add(values.Day)
	case Weekly:
		return s.StartedAt.AddDate(0, 0, 7)
	case Monthly:
		return s.StartedAt.AddDate(0, 1, 0)
	default:
		return values.EmptyTime
	}
}
