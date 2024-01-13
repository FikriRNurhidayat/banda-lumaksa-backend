package subscription

import (
	"context"

	"github.com/google/uuid"
)

type GetSubscriptionUseCase interface {
	Call(ctx context.Context, params *GetSubscriptionParams) (*GetSubscriptionResult, error)
}
type GetSubscriptionParams struct {
	ID uuid.UUID
}
type GetSubscriptionResult struct {
	Subscription Subscription
}
type GetSubscriptionUseCaseImpl struct {
	subscriptionRepository Repository
}

func (u *GetSubscriptionUseCaseImpl) Call(ctx context.Context, params *GetSubscriptionParams) (*GetSubscriptionResult, error) {
	s, err := u.subscriptionRepository.Get(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &GetSubscriptionResult{
		Subscription: s,
	}, nil
}

func NewGetSubscriptionUseCase() GetSubscriptionUseCase {
	return &GetSubscriptionUseCaseImpl{
		subscriptionRepository: nil,
	}
}
