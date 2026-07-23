package application

import (
	"context"

	"vault/src/features/addresses/domain/dto/request"
	"vault/src/features/addresses/domain/dto/response"
	"vault/src/features/addresses/domain/entities"
	"vault/src/features/addresses/domain/repositories"
)

type CreateAddressUseCase struct {
	repo repositories.AddressRepository
}

func NewCreateAddressUseCase(repo repositories.AddressRepository) *CreateAddressUseCase {
	return &CreateAddressUseCase{repo: repo}
}

func (uc *CreateAddressUseCase) Execute(ctx context.Context, userID string, req request.CreateAddressRequest) (response.AddressResponse, error) {
	if err := req.Validate(); err != nil {
		return response.AddressResponse{}, err
	}

	created, err := uc.repo.Create(ctx, entities.Address{
		UserID:         userID,
		Label:          req.Label,
		Recipient:      req.Recipient,
		Phone:          req.Phone,
		Street:         req.Street,
		City:           req.City,
		State:          req.State,
		PostalCode:     req.PostalCode,
		ReferenceNotes: req.References,
	})
	if err != nil {
		return response.AddressResponse{}, err
	}

	return response.FromEntity(created), nil
}
