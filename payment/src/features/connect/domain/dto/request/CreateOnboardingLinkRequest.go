package request

// CreateOnboardingLinkRequest -- igual que en subscriptions/, payment/ no
// tiene su propia base de datos todavía y no puede resolver el email del
// usuario a partir del JWT, así que el cliente lo manda explícito.
type CreateOnboardingLinkRequest struct {
	Email      string `json:"email"`
	RefreshURL string `json:"refresh_url"`
	ReturnURL  string `json:"return_url"`
}
