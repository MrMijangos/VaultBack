package repositories

import (
	"context"
	"errors"
)

var ErrBusinessNotFound = errors.New("el negocio no existe")
var ErrNotOwner = errors.New("no eres el dueno de este negocio")

// BusinessOwnerProvider resuelve el dueño de un negocio para autorizar
// creación/edición/borrado de sus servicios -- vive en este dominio (no en
// businesses) para no acoplar el paquete, mismo patrón que RatingProvider en
// businesses/domain/repositories.
type BusinessOwnerProvider interface {
	GetOwnerUserID(ctx context.Context, businessID string) (string, error)
}
