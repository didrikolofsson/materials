// Package customerrors provides custom error constants for the application.
package customerrors

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInternal           = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrTooManyRequests    = errors.New("too many requests")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrGatewayTimeout     = errors.New("gateway timeout")
)
