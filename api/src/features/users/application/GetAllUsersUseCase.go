package application

import (
	"context"

	"vault/src/features/users/domain/dto/response"
	"vault/src/features/users/domain/repositories"
)

type GetAllUsersUseCase struct {
	repo repositories.UserRepository
}

func NewGetAllUsersUseCase(repo repositories.UserRepository) *GetAllUsersUseCase {
	return &GetAllUsersUseCase{repo: repo}
}

func (uc *GetAllUsersUseCase) Execute(ctx context.Context) ([]response.UserResponse, error) {
	users, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(users), nil
}
