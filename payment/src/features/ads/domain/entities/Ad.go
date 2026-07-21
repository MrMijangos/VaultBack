package entities

import "time"

const (
	AdStatusActive   = "active"
	AdStatusInactive = "inactive"
)

const (
	SectionMarketplace = "marketplace"
	SectionFeed        = "feed"
)

var ValidSections = []string{SectionMarketplace, SectionFeed}

// Ad es un anuncio pagado que un vendedor/restaurador/servicio publica
// mientras tiene una suscripción activa. Se desactiva (no se borra) cuando
// la suscripción se cancela -- ver CancelSubscriptionUseCase -- para poder
// reactivarlo si vuelve a suscribirse sin perder el histórico de
// impressions/clicks.
type Ad struct {
	ID             string
	UserID         string
	SubscriptionID string
	Title          string
	Description    string
	ImageURL       string
	TargetSection  string // "marketplace" o "feed"
	TargetID       string // ID del producto/publicación que promociona
	Status         string
	Impressions    int64
	Clicks         int64
	CreatedAt      time.Time
}

func (a *Ad) IsActive() bool {
	return a.Status == AdStatusActive
}
