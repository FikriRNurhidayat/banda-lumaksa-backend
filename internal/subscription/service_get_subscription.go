package subscription

import (
	"context"

	"github.com/google/uuid"
)

type GetSubscriptionService interface {
	Call(ctx context.Context, params *GetSubscriptionParams) (*GetSubscriptionResult, error)
}

type GetSubscriptionParams struct {
	ID uuid.UUID
}

type GetSubscriptionResult struct {
	Subscription Subscription
}

type GetSubscriptionServiceImpl struct {
	subscriptionRepository Repository
}

func (u *GetSubscriptionServiceImpl) Call(ctx context.Context, params *GetSubscriptionParams) (*GetSubscriptionResult, error) {
	s, err := u.subscriptionRepository.Get(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &GetSubscriptionResult{
		Subscription: s,
	}, nil
}

func NewGetSubscriptionService(subscriptionRepository Repository) GetSubscriptionService {
	return &GetSubscriptionServiceImpl{
		subscriptionRepository: subscriptionRepository,
	}
}
