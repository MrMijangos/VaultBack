package application

import (
	"context"

	"vault/src/features/fcmtokens/domain/dto/request"
	"vault/src/features/fcmtokens/domain/repositories"
)

type RegisterFCMTokenUseCase struct {
	repo repositories.FCMTokenRepository
}

func NewRegisterFCMTokenUseCase(repo repositories.FCMTokenRepository) *RegisterFCMTokenUseCase {
	return &RegisterFCMTokenUseCase{repo: repo}
}

func (uc *RegisterFCMTokenUseCase) Execute(ctx context.Context, userID string, req request.RegisterFCMTokenRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.repo.Register(ctx, userID, req.Token, req.Platform)
}
