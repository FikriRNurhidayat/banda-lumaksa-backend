package subscription_repository

import (
	common_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/entity"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/specification"
)

type SubscriptionRepository common_repository.Repository[subscription_entity.Subscription, subscription_specification.SubscriptionSpecification]
