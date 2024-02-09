package subscription_service

import (
	"context"

	common_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/specification"
	"github.com/google/uuid"
)

type CancelSubscriptionParams struct {
	ID uuid.UUID
}

type CancelSubscriptionResult struct{}

func (s *SubscriptionServiceImpl) CancelSubscription(ctx context.Context, params *CancelSubscriptionParams) (*CancelSubscriptionResult, error) {
	subscription, err := s.subscriptionRepository.Get(ctx, subscription_specification.WithID(params.ID))
	if err != nil {
		return nil, err
	}

	if err := s.subscriptionRepository.Delete(ctx, subscription_specification.WithID(subscription.ID)); err != nil {
		return nil, common_errors.ErrInternalServer
	}

	return &CancelSubscriptionResult{}, nil
}

