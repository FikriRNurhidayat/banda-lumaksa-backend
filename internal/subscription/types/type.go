package subscription_types

import "encoding/json"

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
