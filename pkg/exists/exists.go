package exists

import (
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/values"
)

func String(value string) bool {
	return value != ""
}

func Date(value time.Time) bool {
	return value != values.EmptyTime
}
