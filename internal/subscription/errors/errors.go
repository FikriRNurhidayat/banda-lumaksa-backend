package subscription_errors

import (
	"net/http"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
)

var (
	ErrSubscriptionNotFound = &common_errors.Error{
		Code:    http.StatusNotFound,
		Reason:  "SUBSCRIPTION_NOT_FOUND_ERROR",
		Message: "Subscription not found. Please pass valid subscription id.",
	}

	ErrSubscriptionPastDueAt = &common_errors.Error{
		Code:    http.StatusUnprocessableEntity,
		Reason:  "SUBSCRIPTION_PAST_DUE_AT_ERROR",
		Message: "Due at is in the past. Please pass due at that is on the future.",
	}

	ErrSubscriptionAlreadyExist = &common_errors.Error{
		Code:    http.StatusUnprocessableEntity,
		Reason:  "SUBSCRIPTION_ALREADY_EXIST_ERROR",
		Message: "Subscription already exists. Please use different name.",
	}
)
