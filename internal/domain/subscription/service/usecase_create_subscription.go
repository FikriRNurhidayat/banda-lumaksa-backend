package subscription_service

import (
	"context"
	"time"

	"github.com/google/uuid"

	common_values "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/values"
	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/entity"
	subscription_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/errors"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/specification"
	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/types"
)

type CreateSubscriptionParams struct {
	Name      string
	Fee       int32
	Type      subscription_types.Type
	StartedAt time.Time
	EndedAt   time.Time
	DueAt     time.Time
}

type CreateSubscriptionResult struct {
	Subscription subscription_entity.Subscription
}

func (s *SubscriptionServiceImpl) CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (*CreateSubscriptionResult, error) {
	now := time.Now()
	subscription := subscription_entity.Subscription{
		ID:        uuid.New(),
		Name:      params.Name,
		Fee:       params.Fee,
		Type:      params.Type,
		StartedAt: params.StartedAt,
		EndedAt:   params.EndedAt,
		DueAt:     params.DueAt,
	}

	if subscription.Type == subscription_types.NoType {
		return nil, subscription_errors.ErrSubscriptionTypeInvalid
	}

	if subscription.DueAt == common_values.NoTime {
		subscription.DueAt = s.computeDueAt(subscription, params.StartedAt)
	}

	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	exist, err := s.subscriptionRepository.Exist(ctx, subscription_specification.NameIs(subscription.Name))
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, subscription_errors.ErrSubscriptionAlreadyExist
	}

	if err := s.subscriptionRepository.Save(ctx, subscription); err != nil {
		return nil, err
	}

	return &CreateSubscriptionResult{
		Subscription: subscription,
	}, nil}

