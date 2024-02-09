package subscription_service

import (
	"context"

	"github.com/google/uuid"

	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/entity"
	subscription_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/errors"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/specification"
)

type GetSubscriptionParams struct {
	ID uuid.UUID
}

type GetSubscriptionResult struct {
	Subscription subscription_entity.Subscription
}

func (s *SubscriptionServiceImpl) GetSubscription(ctx context.Context, params *GetSubscriptionParams) (*GetSubscriptionResult, error) {
	subscription, err := s.subscriptionRepository.Get(ctx, subscription_specification.WithID(params.ID))
	if err != nil {
		return nil, err
	}

	if subscription == subscription_entity.NoSubscription {
		return nil, subscription_errors.ErrSubscriptionNotFound
	}

	return &GetSubscriptionResult{
		Subscription: subscription,
	}, nil
}
