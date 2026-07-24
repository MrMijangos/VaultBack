package request

import "errors"

type AddRolesRequest struct {
	Roles []string `json:"roles"`
}

func (r AddRolesRequest) Validate() error {
	if len(r.Roles) == 0 {
		return errors.New("debes indicar al menos un rol")
	}
	for _, role := range r.Roles {
		if !allowedRoles[role] {
			return errors.New("el rol no es valido")
		}
	}
	return nil
}
