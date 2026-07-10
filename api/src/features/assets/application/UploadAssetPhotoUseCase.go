package application

import (
	"context"
	"io"

	"vault/src/core/cloudinary"
	"vault/src/features/assets/domain/dto/response"
	"vault/src/features/assets/domain/repositories"
)

type UploadAssetPhotoUseCase struct {
	repo     repositories.AssetRepository
	uploader *cloudinary.ImageUploader
}

func NewUploadAssetPhotoUseCase(repo repositories.AssetRepository, uploader *cloudinary.ImageUploader) *UploadAssetPhotoUseCase {
	return &UploadAssetPhotoUseCase{repo: repo, uploader: uploader}
}

func (uc *UploadAssetPhotoUseCase) Execute(ctx context.Context, assetID string, userID string, file io.Reader) (response.AssetResponse, error) {
	asset, err := uc.repo.FindByID(ctx, assetID)
	if err != nil {
		return response.AssetResponse{}, err
	}
	if asset.UserID != userID {
		return response.AssetResponse{}, repositories.ErrAssetNotFound
	}

	url, err := uc.uploader.Upload(ctx, file, "vault/assets")
	if err != nil {
		return response.AssetResponse{}, err
	}

	if _, err := uc.repo.AddPhoto(ctx, assetID, url); err != nil {
		return response.AssetResponse{}, err
	}

	photos, err := uc.repo.FindPhotosByAssetID(ctx, assetID)
	if err != nil {
		return response.AssetResponse{}, err
	}

	return response.FromEntity(asset, photos), nil
}
