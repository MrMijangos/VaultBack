package controllers

import (
	"errors"
	"net/http"

	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/orders/application"
)

func statusForError(err error) int {
	switch {
	case errors.Is(err, application.ErrInvalidRequest):
		return http.StatusBadRequest
	case errors.Is(err, application.ErrSellerNotOnboarded):
		return http.StatusConflict
	case errors.Is(err, application.ErrOrderNotFound):
		return http.StatusNotFound
	case errors.Is(err, application.ErrNotBuyer):
		return http.StatusForbidden
	case errors.Is(err, application.ErrNotHeld):
		return http.StatusConflict
	case stripeclient.IsNotConfigured(err):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
