package application

import (
	"context"

	assetrepositories "vault/src/features/assets/domain/repositories"
	"vault/src/features/blockchaincertificates/domain/dto/request"
	"vault/src/features/blockchaincertificates/domain/dto/response"
	"vault/src/features/blockchaincertificates/domain/entities"
	"vault/src/features/blockchaincertificates/domain/repositories"
)

type CreateBlockchainCertificateUseCase struct {
	repo      repositories.BlockchainCertificateRepository
	assetRepo assetrepositories.AssetRepository
}

func NewCreateBlockchainCertificateUseCase(repo repositories.BlockchainCertificateRepository, assetRepo assetrepositories.AssetRepository) *CreateBlockchainCertificateUseCase {
	return &CreateBlockchainCertificateUseCase{repo: repo, assetRepo: assetRepo}
}

func (uc *CreateBlockchainCertificateUseCase) Execute(ctx context.Context, userID string, req request.CreateBlockchainCertificateRequest) (response.BlockchainCertificateResponse, error) {
	if err := req.Validate(); err != nil {
		return response.BlockchainCertificateResponse{}, err
	}

	asset, err := uc.assetRepo.FindByID(ctx, req.AssetID)
	if err != nil {
		return response.BlockchainCertificateResponse{}, err
	}
	if asset.UserID != userID {
		return response.BlockchainCertificateResponse{}, assetrepositories.ErrAssetNotFound
	}

	taken, err := uc.repo.ExistsByTxID(ctx, req.TxID)
	if err != nil {
		return response.BlockchainCertificateResponse{}, err
	}
	if taken {
		return response.BlockchainCertificateResponse{}, repositories.ErrTxIDAlreadyExists
	}

	created, err := uc.repo.Create(ctx, entities.BlockchainCertificate{
		AssetID:   req.AssetID,
		OwnerID:   userID,
		TxID:      req.TxID,
		AssetHash: req.AssetHash,
		Action:    req.Action,
		Network:   req.Network,
	})
	if err != nil {
		return response.BlockchainCertificateResponse{}, err
	}

	return response.FromEntity(created), nil
}
