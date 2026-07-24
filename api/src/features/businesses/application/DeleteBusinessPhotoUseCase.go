package application

import (
	"context"

	"vault/src/features/businesses/domain/dto/response"
	"vault/src/features/businesses/domain/repositories"
)

type DeleteBusinessPhotoUseCase struct {
	repo repositories.BusinessRepository
}

func NewDeleteBusinessPhotoUseCase(repo repositories.BusinessRepository) *DeleteBusinessPhotoUseCase {
	return &DeleteBusinessPhotoUseCase{repo: repo}
}

func (uc *DeleteBusinessPhotoUseCase) Execute(ctx context.Context, businessID string, photoID string, userID string) (response.BusinessResponse, error) {
	business, err := uc.repo.FindByID(ctx, businessID)
	if err != nil {
		return response.BusinessResponse{}, err
	}
	if business.UserID != userID {
		return response.BusinessResponse{}, repositories.ErrBusinessNotFound
	}

	if err := uc.repo.DeletePhoto(ctx, photoID, businessID); err != nil {
		return response.BusinessResponse{}, err
	}

	photos, err := uc.repo.FindPhotosByBusinessID(ctx, businessID)
	if err != nil {
		return response.BusinessResponse{}, err
	}

	return response.FromEntity(business, photos), nil
}
