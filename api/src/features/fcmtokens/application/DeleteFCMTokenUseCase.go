package application

import (
	"context"

	"vault/src/features/fcmtokens/domain/dto/request"
	"vault/src/features/fcmtokens/domain/repositories"
)

type DeleteFCMTokenUseCase struct {
	repo repositories.FCMTokenRepository
}

func NewDeleteFCMTokenUseCase(repo repositories.FCMTokenRepository) *DeleteFCMTokenUseCase {
	return &DeleteFCMTokenUseCase{repo: repo}
}

func (uc *DeleteFCMTokenUseCase) Execute(ctx context.Context, req request.DeleteFCMTokenRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	// No se valida que el token pertenezca a userID -- al hacer logout la
	// app solo conoce el token del propio dispositivo, y borrar un token que
	// ya no existe (o que es de otro usuario) es un no-op inofensivo, no un
	// error que valga la pena reportar.
	return uc.repo.DeleteToken(ctx, req.Token)
}
