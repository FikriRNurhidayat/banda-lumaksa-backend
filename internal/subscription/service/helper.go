package subscription_service

import (
	"context"
	"fmt"
	"time"

	common_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	common_values "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/values"
	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/entity"
	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/types"
	
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/transaction"
	"github.com/google/uuid"
)

func (s *SubscriptionServiceImpl) computeDueAt(subscription subscription_entity.Subscription, startFrom time.Time) time.Time {
	switch subscription.Type {
	case subscription_types.Daily:
		return startFrom.Add(common_values.Day)
	case subscription_types.Weekly:
		return startFrom.AddDate(0, 0, 7)
	case subscription_types.Monthly:
		return startFrom.AddDate(0, 1, 0)
	default:
		return common_values.NoTime
	}
}

func (s *SubscriptionServiceImpl) chargeSubscription(ctx context.Context, subscription subscription_entity.Subscription) (subscription_entity.Subscription, error) {
	now := time.Now()
	subscription.UpdatedAt = now
	subscription.DueAt = s.computeDueAt(subscription, now)

	transaction := transaction.Transaction{
		ID:          uuid.New(),
		Description: subscription.GetTransactionDescription(),
		Amount:      subscription.Fee,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.transactionManager.Execute(ctx, func(ctx context.Context) error {
		if err := s.subscriptionRepository.Save(ctx, subscription); err != nil {
			return err
		}

		if err := s.transactionRepository.Save(ctx, transaction); err != nil {
			return err
		}

		return nil
	}); err != nil {
		fmt.Println(err.Error())
		return subscription_entity.NoSubscription, common_errors.ErrInternalServer
	}

	return subscription, nil
}
