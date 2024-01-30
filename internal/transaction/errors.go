package transaction

import (
	"net/http"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
)

var (
	ErrTransactionNotFound = &errors.Error{
		Code:    http.StatusNotFound,
		Reason:  "TRANSACTION_NOT_FOUND_ERROR",
		Message: "Transaction not found. Please pass valid transaction id.",
	}
)
