package request

import (
	"errors"

	"vault/src/core/validation"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if r.Email == "" || r.Password == "" {
		return errors.New("correo y contraseña son obligatorios")
	}
	if !validation.IsValidEmail(r.Email) {
		return errors.New("el correo no es valido")
	}
	return nil
}
