package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/fcmtokens/application"
	"vault/src/features/fcmtokens/infrastructure/adapters"
	"vault/src/features/fcmtokens/infrastructure/controllers"
)

func BuildRegisterFCMTokenController(pool *pgxpool.Pool) *controllers.RegisterFCMTokenController {
	repo := adapters.NewPostgreSQLFCMTokenRepository(pool)
	useCase := application.NewRegisterFCMTokenUseCase(repo)
	return controllers.NewRegisterFCMTokenController(useCase)
}

func BuildDeleteFCMTokenController(pool *pgxpool.Pool) *controllers.DeleteFCMTokenController {
	repo := adapters.NewPostgreSQLFCMTokenRepository(pool)
	useCase := application.NewDeleteFCMTokenUseCase(repo)
	return controllers.NewDeleteFCMTokenController(useCase)
}

// BuildFCMTokenProvider expone el adapter crudo para quien necesite el
// puerto push.TokenProvider (main.go, al armar el push.Sender) -- lo
// satisface por estructura, sin que src/core/push dependa de esta feature.
func BuildFCMTokenProvider(pool *pgxpool.Pool) *adapters.PostgreSQLFCMTokenRepository {
	return adapters.NewPostgreSQLFCMTokenRepository(pool)
}
