package request

import (
	"errors"
	"strings"
)

var allowedRoles = map[string]bool{
	"usuario":     true,
	"vendedor":    true,
	"restaurador": true,
	"servicio":    true,
	"admin":       true,
}

type CreateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
}

func (r *CreateUserRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("el nombre es obligatorio")
	}
	if !strings.Contains(r.Email, "@") {
		return errors.New("el correo no es valido")
	}
	if len(r.Password) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}
	if r.Role == "" {
		r.Role = "usuario"
	}
	if !allowedRoles[r.Role] {
		return errors.New("el rol no es valido")
	}
	return nil
}
