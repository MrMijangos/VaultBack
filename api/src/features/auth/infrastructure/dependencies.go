package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/auth/application"
	"vault/src/features/auth/infrastructure/adapters"
	"vault/src/features/auth/infrastructure/controllers"
)

func BuildLoginController(pool *pgxpool.Pool, jwtSecret string, cookieSecure bool) *controllers.LoginController {
	repo := adapters.NewPostgreSQLAuthRepository(pool)
	useCase := application.NewLoginUseCase(repo, jwtSecret)
	return controllers.NewLoginController(useCase, cookieSecure)
}
