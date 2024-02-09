package subscription_service

import (
	"context"

	common_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/entity"
	subscription_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/errors"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/specification"
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
		return nil, common_errors.ErrInternalServer
	}

	return &ChargeSubscriptionResult{}, nil
}
