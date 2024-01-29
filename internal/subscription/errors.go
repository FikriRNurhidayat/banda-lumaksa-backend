package subscription

import (
	"net/http"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/errors"
)

var (
	ErrSubscriptionNotFound = &errors.Error{
		Code:    http.StatusNotFound,
		Reason:  "SUBSCRIPTION_NOT_FOUND_ERROR",
		Message: "Subscription not found. Please pass valid subscription id.",
	}

	ErrSubscriptionPastDueAt = &errors.Error{
		Code:    http.StatusUnprocessableEntity,
		Reason:  "SUBSCRIPTION_PAST_DUE_AT_ERROR",
		Message: "Due at is in the past. Please pass due at that is on the future.",
	}
)
