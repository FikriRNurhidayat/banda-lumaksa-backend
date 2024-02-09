package subscription_service

import (
	"context"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/transaction"

	common_manager "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"

	subscription_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/repository"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (*CreateSubscriptionResult, error)
	GetSubscription(ctx context.Context, params *GetSubscriptionParams) (*GetSubscriptionResult, error)
	ListSubscriptions(ctx context.Context, params *ListSubscriptionsParams) (*ListSubscriptionsResult, error)
	CancelSubscription(ctx context.Context, params *CancelSubscriptionParams) (*CancelSubscriptionResult, error)
	ChargeSubscription(ctx context.Context, params *ChargeSubscriptionParams) (*ChargeSubscriptionResult, error)
	ChargeSubscriptions(ctx context.Context, params *ChargeSubscriptionsParams) (*ChargeSubscriptionsResult, error)
}

type SubscriptionServiceImpl struct {
	subscriptionRepository subscription_repository.SubscriptionRepository
	transactionRepository  transaction.TransactionRepository
	transactionManager     common_manager.TransactionManager
	logger                 logger.Logger
}

func New(
	logger logger.Logger,
	subscriptionRepository subscription_repository.SubscriptionRepository,
	transactionRepository transaction.TransactionRepository,
	transactionManager common_manager.TransactionManager) SubscriptionService {
	return &SubscriptionServiceImpl{
		subscriptionRepository: subscriptionRepository,
		transactionRepository:  transactionRepository,
		transactionManager:     transactionManager,
		logger:                 logger,
	}
}
