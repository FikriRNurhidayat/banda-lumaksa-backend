package subscription

import (
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
)

type SubscriptionRepository repository.Repository[Subscription, SubscriptionSpecification]
