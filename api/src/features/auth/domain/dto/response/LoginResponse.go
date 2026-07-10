package response

import "vault/src/features/auth/domain/entities"

type LoginResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
}

func FromCredentials(c entities.Credentials) LoginResponse {
	return LoginResponse{
		ID:        c.UserID,
		Name:      c.Name,
		Email:     c.Email,
		AvatarURL: c.AvatarURL,
		Role:      c.Role,
	}
}
