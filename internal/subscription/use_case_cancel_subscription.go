package subscription

import (
	"context"

	"github.com/google/uuid"
)

type CancelSubscriptionUseCase interface {
	Call(ctx context.Context, params *CancelSubscriptionParams) (*CancelSubscriptionResult, error)
}
type CancelSubscriptionParams struct {
	ID uuid.UUID
}
type CancelSubscriptionResult struct{}
type CancelSubscriptionUseCaseImpl struct {
	subscriptionRepository Repository
}

func (u *CancelSubscriptionUseCaseImpl) Call(ctx context.Context, params *CancelSubscriptionParams) (*CancelSubscriptionResult, error) {
	s, err := u.subscriptionRepository.Get(ctx, params.ID)
	if err != nil {
		return nil, err
	}

  if err := u.subscriptionRepository.Delete(ctx, s.ID); err != nil {
    return nil, err
  }

  return &CancelSubscriptionResult{}, nil
}

func NewCancelSubscriptionUseCase() CancelSubscriptionUseCase {
	return &CancelSubscriptionUseCaseImpl{
		subscriptionRepository: nil,
	}
}
