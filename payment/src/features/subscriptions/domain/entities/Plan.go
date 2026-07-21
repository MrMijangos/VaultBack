package entities

// Plan representa uno de los 3 planes fijos de suscripción (básico/pro/premium).
// No se guardan en base de datos todavía -- ver PlanRepository en memoria --
// porque el esquema de Supabase para esto se deja para el final.
type Plan struct {
	ID            string
	Name          string
	PriceMXN      float64
	StripePriceID string
	MaxAds        int
	// TargetSections son las secciones donde puede aparecer un anuncio de
	// este plan. "marketplace" siempre está permitido; "feed" solo en pro y
	// premium, y ahí el vendedor elige en cuál publicar cada anuncio.
	TargetSections []string
	// CommissionRate es el % (como fracción, ej. 0.08 = 8%) que Vault retiene
	// de cada venta liberada del escrow para vendedores en este plan.
	CommissionRate float64
}

// DefaultCommissionRate se usa para vendedores sin una suscripción activa --
// no bloqueamos la venta por no estar suscritos, pero no reciben el
// descuento de comisión de los planes pagados (aplica la misma tasa que
// básico, el plan de entrada).
const DefaultCommissionRate = 0.08

const (
	SectionMarketplace = "marketplace"
	SectionFeed        = "feed"
)

const (
	PlanBasico  = "basico"
	PlanPro     = "pro"
	PlanPremium = "premium"
)
