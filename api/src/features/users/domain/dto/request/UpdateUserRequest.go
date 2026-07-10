package request

import (
	"errors"
	"strings"
)

type UpdateUserRequest struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
}

func (r UpdateUserRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("el nombre es obligatorio")
	}
	if r.Role != "" && !allowedRoles[r.Role] {
		return errors.New("el rol no es valido")
	}
	return nil
}
