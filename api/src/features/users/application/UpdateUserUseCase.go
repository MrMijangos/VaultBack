package application

import (
	"context"

	"vault/src/features/users/domain/dto/request"
	"vault/src/features/users/domain/dto/response"
	"vault/src/features/users/domain/repositories"
)

type UpdateUserUseCase struct {
	repo repositories.UserRepository
}

func NewUpdateUserUseCase(repo repositories.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, id string, req request.UpdateUserRequest) (response.UserResponse, error) {
	if err := req.Validate(); err != nil {
		return response.UserResponse{}, err
	}

	current, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.UserResponse{}, err
	}

	current.Name = req.Name
	current.AvatarURL = req.AvatarURL
	if req.Role != "" {
		current.Role = req.Role
	}

	updated, err := uc.repo.Update(ctx, id, current)
	if err != nil {
		return response.UserResponse{}, err
	}

	return response.FromEntity(updated), nil
}
