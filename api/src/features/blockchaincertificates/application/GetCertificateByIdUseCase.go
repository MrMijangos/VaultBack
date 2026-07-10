package application

import (
	"context"

	"vault/src/features/blockchaincertificates/domain/dto/response"
	"vault/src/features/blockchaincertificates/domain/repositories"
)

type GetCertificateByIdUseCase struct {
	repo repositories.BlockchainCertificateRepository
}

func NewGetCertificateByIdUseCase(repo repositories.BlockchainCertificateRepository) *GetCertificateByIdUseCase {
	return &GetCertificateByIdUseCase{repo: repo}
}

func (uc *GetCertificateByIdUseCase) Execute(ctx context.Context, id string) (response.BlockchainCertificateResponse, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.BlockchainCertificateResponse{}, err
	}
	return response.FromEntity(c), nil
}
