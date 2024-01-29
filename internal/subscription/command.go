package subscription

import "context"

type SubscriptionCommand interface {
	ChargeSubscriptions(ctx context.Context) error
}

type SubscriptionCommandImpl struct {
	subscriptionService SubscriptionService
}

func (c *SubscriptionCommandImpl) ChargeSubscriptions(ctx context.Context) error {
	_, err := c.subscriptionService.ChargeSubscriptions(ctx, &ChargeSubscriptionsParams{})
	return err
}

func NewSubscriptionCommand(subscriptionService SubscriptionService) SubscriptionCommand {
	return &SubscriptionCommandImpl{
		subscriptionService: subscriptionService,
	}
}
