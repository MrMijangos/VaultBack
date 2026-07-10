package application

import (
	"context"

	"vault/src/features/users/domain/dto/response"
	"vault/src/features/users/domain/repositories"
)

type GetUserByIdUseCase struct {
	repo repositories.UserRepository
}

func NewGetUserByIdUseCase(repo repositories.UserRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{repo: repo}
}

func (uc *GetUserByIdUseCase) Execute(ctx context.Context, id string) (response.UserResponse, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.UserResponse{}, err
	}
	return response.FromEntity(user), nil
}
