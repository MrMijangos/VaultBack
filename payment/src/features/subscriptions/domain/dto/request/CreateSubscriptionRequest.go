package request

// CreateSubscriptionRequest -- payment/ no tiene su propia base de datos
// todavía, así que no puede resolver el email del usuario a partir del JWT
// (que solo trae user_id/role); el cliente (Flutter) lo manda explícito.
// PaymentMethodID lo genera flutter_stripe en el dispositivo -- este backend
// nunca ve el número de tarjeta.
type CreateSubscriptionRequest struct {
	PlanID          string `json:"plan_id"`
	Email           string `json:"email"`
	PaymentMethodID string `json:"payment_method_id"`
}
