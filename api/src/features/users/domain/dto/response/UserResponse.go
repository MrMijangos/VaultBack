package response

import "vault/src/features/users/domain/entities"

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
}

func FromEntity(user entities.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
	}
}

func FromEntities(users []entities.User) []UserResponse {
	list := make([]UserResponse, 0, len(users))
	for _, u := range users {
		list = append(list, FromEntity(u))
	}
	return list
}
