package subscription_entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/types"
)

type Subscription struct {
	ID        uuid.UUID
	Name      string
	Fee       int32
	Type      subscription_types.Type
	StartedAt time.Time
	EndedAt   time.Time
	DueAt     time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Subscriptions []Subscription

var NoSubscription = Subscription{}
var NoSubscriptions = []Subscription{}

func (s Subscription) GetTransactionDescription() string {
	return fmt.Sprintf("Pembayaran biaya langganan untuk layanan %s, senilai %d.", s.Name, s.Fee)
}
