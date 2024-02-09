package subscription_command

import (
	"context"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/service"
)

type SubscriptionCommand interface {
	ChargeSubscriptions(ctx context.Context) error
}

type SubscriptionCommandImpl struct {
	subscriptionService subscription_service.SubscriptionService
}

func (c *SubscriptionCommandImpl) ChargeSubscriptions(ctx context.Context) error {
	_, err := c.subscriptionService.ChargeSubscriptions(ctx, &subscription_service.ChargeSubscriptionsParams{})
	return err
}

func New(subscriptionService subscription_service.SubscriptionService) SubscriptionCommand {
	return &SubscriptionCommandImpl{
		subscriptionService: subscriptionService,
	}
}
