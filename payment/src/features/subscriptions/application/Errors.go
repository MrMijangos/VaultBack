package application

import "errors"

var (
	ErrRoleNotAllowed    = errors.New("tu tipo de cuenta no puede suscribirse (solo vendedor, restaurador o servicio)")
	ErrAlreadySubscribed = errors.New("ya tienes una suscripción activa")
	ErrNotSubscribed     = errors.New("no tienes una suscripción activa")
	ErrInvalidRequest    = errors.New("plan_id, email y payment_method_id son obligatorios")
)

// allowedRoles son los tipos de cuenta que pueden suscribirse -- deben
// coincidir con los roles que emite api/ en el JWT (Claims.Role).
var allowedRoles = map[string]bool{
	"vendedor":    true,
	"restaurador": true,
	"servicio":    true,
}

func isRoleAllowed(role string) bool {
	return allowedRoles[role]
}
