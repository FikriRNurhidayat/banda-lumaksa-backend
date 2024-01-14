package subscription

import (
	"time"

	"github.com/google/uuid"
)

type Type int

const (
	Daily Type = iota
	Weekly
	Monthly
	Yearly
	InvalidType
)

func (t Type) String() string {
	switch t {
	case Daily:
		return "Daily"
	case Weekly:
		return "Weekly"
	case Monthly:
		return "Monthly"
	case Yearly:
		return "Yearly"
	default:
		return ""
	}
}

func GetType(str string) Type {
	switch str {
	case "Daily":
		return Daily
	case "Weekly":
		return Weekly
	case "Monthly":
		return Monthly
	case "Yearly":
		return Yearly
	default:
		return -1
	}
}

type Subscription struct {
	ID        uuid.UUID
	Name      string
	Fee       int32
	Type      Type
	StartedAt time.Time
	EndedAt   time.Time
	DueAt     time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Subscriptions []Subscription

var EmptySubscription = Subscription{}
var EmptySubscriptions = []Subscription{}
