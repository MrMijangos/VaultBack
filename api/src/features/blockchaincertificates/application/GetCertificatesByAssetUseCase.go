package application

import (
	"context"

	"vault/src/features/blockchaincertificates/domain/dto/response"
	"vault/src/features/blockchaincertificates/domain/repositories"
)

type GetCertificatesByAssetUseCase struct {
	repo repositories.BlockchainCertificateRepository
}

func NewGetCertificatesByAssetUseCase(repo repositories.BlockchainCertificateRepository) *GetCertificatesByAssetUseCase {
	return &GetCertificatesByAssetUseCase{repo: repo}
}

func (uc *GetCertificatesByAssetUseCase) Execute(ctx context.Context, assetID string) ([]response.BlockchainCertificateResponse, error) {
	list, err := uc.repo.FindByAssetID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
