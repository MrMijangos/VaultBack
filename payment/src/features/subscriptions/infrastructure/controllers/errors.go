package controllers

import (
	"errors"
	"net/http"

	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/subscriptions/application"
)

// statusForError traduce errores de negocio a códigos HTTP. Cualquier error
// no reconocido (por ejemplo, uno que venga de Stripe) se responde como 500
// pero con el mensaje real -- este servicio no tiene usuarios finales
// directos, solo Flutter, así que no hace falta ocultar el detalle.
func statusForError(err error) int {
	switch {
	case errors.Is(err, application.ErrRoleNotAllowed):
		return http.StatusForbidden
	case errors.Is(err, application.ErrAlreadySubscribed):
		return http.StatusConflict
	case errors.Is(err, application.ErrNotSubscribed):
		return http.StatusNotFound
	case errors.Is(err, application.ErrInvalidRequest):
		return http.StatusBadRequest
	case errors.Is(err, stripeclient.ErrNotConfigured):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
