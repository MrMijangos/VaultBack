package request

import "errors"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if r.Email == "" || r.Password == "" {
		return errors.New("correo y contraseña son obligatorios")
	}
	return nil
}
