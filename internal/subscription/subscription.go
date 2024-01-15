package subscription

import (
	"encoding/json"
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

func (t *Type) UnmarshalJSON(b []byte) error {
	var val string
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	*t = GetType(val)
	return nil
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
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

var NoSubscription = Subscription{}
var NoSubscriptions = []Subscription{}
