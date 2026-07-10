package application

import (
	"context"

	"vault/src/features/businesses/domain/repositories"
)

type DeleteBusinessUseCase struct {
	repo repositories.BusinessRepository
}

func NewDeleteBusinessUseCase(repo repositories.BusinessRepository) *DeleteBusinessUseCase {
	return &DeleteBusinessUseCase{repo: repo}
}

func (uc *DeleteBusinessUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
