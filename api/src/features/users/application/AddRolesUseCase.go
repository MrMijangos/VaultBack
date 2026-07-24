package application

import (
	"context"

	"vault/src/features/users/domain/dto/request"
	"vault/src/features/users/domain/dto/response"
	"vault/src/features/users/domain/repositories"
)

type AddRolesUseCase struct {
	repo repositories.UserRepository
}

func NewAddRolesUseCase(repo repositories.UserRepository) *AddRolesUseCase {
	return &AddRolesUseCase{repo: repo}
}

func (uc *AddRolesUseCase) Execute(ctx context.Context, id string, req request.AddRolesRequest) (response.UserResponse, error) {
	if err := req.Validate(); err != nil {
		return response.UserResponse{}, err
	}

	updated, err := uc.repo.AddRoles(ctx, id, req.Roles)
	if err != nil {
		return response.UserResponse{}, err
	}

	return response.FromEntity(updated), nil
}
