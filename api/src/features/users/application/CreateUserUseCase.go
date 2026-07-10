package application

import (
	"context"
	"errors"

	"vault/src/core/security"
	"vault/src/features/users/domain/dto/request"
	"vault/src/features/users/domain/dto/response"
	"vault/src/features/users/domain/entities"
	"vault/src/features/users/domain/repositories"
)

var ErrEmailTaken = errors.New("ya existe una cuenta con ese correo")

type CreateUserUseCase struct {
	repo repositories.UserRepository
}

func NewCreateUserUseCase(repo repositories.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req request.CreateUserRequest) (response.UserResponse, error) {
	if err := req.Validate(); err != nil {
		return response.UserResponse{}, err
	}

	taken, err := uc.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return response.UserResponse{}, err
	}
	if taken {
		return response.UserResponse{}, ErrEmailTaken
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return response.UserResponse{}, err
	}

	user := entities.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		AvatarURL:    req.AvatarURL,
		Role:         req.Role,
	}

	created, err := uc.repo.Create(ctx, user)
	if err != nil {
		return response.UserResponse{}, err
	}

	return response.FromEntity(created), nil
}
