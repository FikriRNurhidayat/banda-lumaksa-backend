package errors

import "net/http"

var (
	ErrInternalServer = &Error{
		Code:    http.StatusInternalServerError,
		Reason:  "INTERNAL_SERVER_ERROR",
		Message: "Internal server error.",
	}
)
