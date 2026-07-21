package application

import "errors"

var (
	ErrNoActiveSubscription = errors.New("necesitas una suscripción activa para publicar anuncios")
	ErrMaxAdsReached        = errors.New("alcanzaste el límite de anuncios activos de tu plan")
	ErrSectionNotAllowed    = errors.New("tu plan no permite anuncios en esa sección")
	ErrAdNotFound           = errors.New("anuncio no encontrado")
	ErrNotOwner             = errors.New("este anuncio no te pertenece")
	ErrInvalidSection       = errors.New("target_section debe ser 'marketplace' o 'feed'")
)
