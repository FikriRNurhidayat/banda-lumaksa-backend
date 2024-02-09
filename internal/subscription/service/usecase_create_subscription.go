package subscription_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	common_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	common_values "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/values"
	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/entity"
	subscription_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/errors"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/specification"
	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/types"
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

	if subscription.DueAt == common_values.NoTime {
		subscription.DueAt = s.computeDueAt(subscription, params.StartedAt)
	}

	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	exist, err := s.subscriptionRepository.Exist(ctx, subscription_specification.NameIs(subscription.Name))
	if err != nil {
		fmt.Println(err.Error())
		return nil, common_errors.ErrInternalServer
	}

	if exist {
		return nil, subscription_errors.ErrSubscriptionAlreadyExist
	}

	if err := s.subscriptionRepository.Save(ctx, subscription); err != nil {
		fmt.Println(err.Error())
		return nil, common_errors.ErrInternalServer
	}

	return &CreateSubscriptionResult{
		Subscription: subscription,
	}, nil}

