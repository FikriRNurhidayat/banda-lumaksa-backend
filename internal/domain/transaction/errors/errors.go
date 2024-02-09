package transaction_errors

import (
	"net/http"

	common_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
)

var (
	ErrTransactionNotFound = &common_errors.Error{
		Code:    http.StatusNotFound,
		Reason:  "TRANSACTION_NOT_FOUND_ERROR",
		Message: "Transaction not found. Please pass valid transaction id.",
	}
)
