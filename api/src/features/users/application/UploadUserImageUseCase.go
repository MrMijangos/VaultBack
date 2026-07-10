package application

import (
	"context"
	"io"

	"vault/src/core/cloudinary"
	"vault/src/features/users/domain/dto/response"
	"vault/src/features/users/domain/repositories"
)

type UploadUserImageUseCase struct {
	repo     repositories.UserRepository
	uploader *cloudinary.ImageUploader
}

func NewUploadUserImageUseCase(repo repositories.UserRepository, uploader *cloudinary.ImageUploader) *UploadUserImageUseCase {
	return &UploadUserImageUseCase{repo: repo, uploader: uploader}
}

func (uc *UploadUserImageUseCase) Execute(ctx context.Context, id string, file io.Reader) (response.UserResponse, error) {
	url, err := uc.uploader.Upload(ctx, file, "vault/users")
	if err != nil {
		return response.UserResponse{}, err
	}

	updated, err := uc.repo.UpdateImage(ctx, id, url)
	if err != nil {
		return response.UserResponse{}, err
	}

	return response.FromEntity(updated), nil
}
