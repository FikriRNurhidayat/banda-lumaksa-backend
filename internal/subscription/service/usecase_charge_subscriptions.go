package subscription_service

import (
	"context"
	"time"

	common_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	common_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/specification"
)

type ChargeSubscriptionsParams struct{}
type ChargeSubscriptionsResult struct{}

func (s *SubscriptionServiceImpl) ChargeSubscriptions(ctx context.Context, params *ChargeSubscriptionsParams) (*ChargeSubscriptionsResult, error) {
	today := time.Now()
	iterator, err := s.subscriptionRepository.Each(ctx, common_repository.ListArgs[subscription_specification.SubscriptionSpecification]{
		Filters: subscription_specification.SubscriptionSpecifications{subscription_specification.DueBefore(today), subscription_specification.NotEnded(today)},
	})
	if err != nil {
		s.logger.Error("subscription repository each", err)
		return nil, common_errors.ErrInternalServer
	}

	for iterator.Next() {
		subscription, err := iterator.Entry()
		if err != nil {
			continue
		}

		if _, err := s.chargeSubscription(ctx, subscription); err != nil {
			continue
		}
	}

	return &ChargeSubscriptionsResult{}, nil
}
