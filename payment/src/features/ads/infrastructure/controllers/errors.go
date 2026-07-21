package controllers

import (
	"errors"
	"net/http"

	"vault-payment/src/features/ads/application"
)

func statusForError(err error) int {
	switch {
	case errors.Is(err, application.ErrNoActiveSubscription):
		return http.StatusForbidden
	case errors.Is(err, application.ErrMaxAdsReached), errors.Is(err, application.ErrSectionNotAllowed):
		return http.StatusConflict
	case errors.Is(err, application.ErrInvalidSection):
		return http.StatusBadRequest
	case errors.Is(err, application.ErrAdNotFound):
		return http.StatusNotFound
	case errors.Is(err, application.ErrNotOwner):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
