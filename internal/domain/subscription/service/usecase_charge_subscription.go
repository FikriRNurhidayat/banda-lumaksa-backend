package subscription_service

import (
	"context"

	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/entity"
	subscription_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/errors"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/specification"
	"github.com/google/uuid"
)

type ChargeSubscriptionParams struct {
	ID uuid.UUID
}

type ChargeSubscriptionResult struct {
	Subscription subscription_entity.Subscription
}

func (s *SubscriptionServiceImpl) ChargeSubscription(ctx context.Context, params *ChargeSubscriptionParams) (*ChargeSubscriptionResult, error) {
	subscription, err := s.subscriptionRepository.Get(ctx, subscription_specification.WithID(params.ID))
	if err != nil {
		return nil, subscription_errors.ErrSubscriptionNotFound
	}

	subscription, err = s.chargeSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return &ChargeSubscriptionResult{}, nil
}
