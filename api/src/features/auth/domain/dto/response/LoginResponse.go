package response

import "vault/src/features/auth/domain/entities"

type LoginResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	AvatarURL string   `json:"avatar_url"`
	Role      string   `json:"role"`
	Roles     []string `json:"roles"`
	Token     string   `json:"token,omitempty"`
	// IsNewUser solo lo llena el login con Google -- el login por correo lo
	// deja siempre en false (la cuenta, si existe, nunca es "nueva" en ese
	// flujo).
	IsNewUser bool `json:"is_new_user"`
}

func FromCredentials(c entities.Credentials) LoginResponse {
	return fromCredentials(c, false)
}

func FromGoogleLogin(c entities.Credentials, isNewUser bool) LoginResponse {
	return fromCredentials(c, isNewUser)
}

func fromCredentials(c entities.Credentials, isNewUser bool) LoginResponse {
	roles := c.Roles
	if roles == nil {
		roles = []string{}
	}
	return LoginResponse{
		ID:        c.UserID,
		Name:      c.Name,
		Email:     c.Email,
		AvatarURL: c.AvatarURL,
		Role:      c.Role,
		Roles:     roles,
		IsNewUser: isNewUser,
	}
}
