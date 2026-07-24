package application

import (
	"context"
	"io"

	"vault/src/core/cloudinary"
	"vault/src/features/businesses/domain/dto/response"
	"vault/src/features/businesses/domain/repositories"
)

type UploadBusinessPhotoUseCase struct {
	repo     repositories.BusinessRepository
	uploader *cloudinary.ImageUploader
}

func NewUploadBusinessPhotoUseCase(repo repositories.BusinessRepository, uploader *cloudinary.ImageUploader) *UploadBusinessPhotoUseCase {
	return &UploadBusinessPhotoUseCase{repo: repo, uploader: uploader}
}

func (uc *UploadBusinessPhotoUseCase) Execute(ctx context.Context, businessID string, userID string, file io.Reader) (response.BusinessResponse, error) {
	business, err := uc.repo.FindByID(ctx, businessID)
	if err != nil {
		return response.BusinessResponse{}, err
	}
	if business.UserID != userID {
		return response.BusinessResponse{}, repositories.ErrBusinessNotFound
	}

	url, err := uc.uploader.Upload(ctx, file, "vault/businesses")
	if err != nil {
		return response.BusinessResponse{}, err
	}

	if _, err := uc.repo.AddPhoto(ctx, businessID, url); err != nil {
		return response.BusinessResponse{}, err
	}

	photos, err := uc.repo.FindPhotosByBusinessID(ctx, businessID)
	if err != nil {
		return response.BusinessResponse{}, err
	}

	return response.FromEntity(business, photos), nil
}
