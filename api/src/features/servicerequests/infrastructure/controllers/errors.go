package controllers

import (
	"errors"
	"net/http"

	assetrepositories "vault/src/features/assets/domain/repositories"
	businessrepositories "vault/src/features/businesses/domain/repositories"
	"vault/src/features/servicerequests/application"
	"vault/src/features/servicerequests/domain/repositories"
)

func statusForError(err error) int {
	switch {
	case errors.Is(err, application.ErrNotOwner):
		return http.StatusForbidden
	case errors.Is(err, application.ErrNotBusinessOwner):
		return http.StatusForbidden
	case errors.Is(err, application.ErrInvalidTransition):
		return http.StatusConflict
	case errors.Is(err, repositories.ErrServiceRequestNotFound):
		return http.StatusNotFound
	case errors.Is(err, assetrepositories.ErrAssetNotFound):
		return http.StatusNotFound
	case errors.Is(err, businessrepositories.ErrBusinessNotFound):
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
