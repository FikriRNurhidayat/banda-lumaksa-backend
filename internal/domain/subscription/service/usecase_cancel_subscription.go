package subscription_service

import (
	"context"

	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/specification"
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
		return nil, err
	}

	return &CancelSubscriptionResult{}, nil
}

