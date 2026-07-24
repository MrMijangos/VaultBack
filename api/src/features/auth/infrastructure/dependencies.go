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

func BuildGoogleLoginController(pool *pgxpool.Pool, googleClientIDs string, jwtSecret string, cookieSecure bool) *controllers.GoogleLoginController {
	repo := adapters.NewPostgreSQLAuthRepository(pool)
	verifier := adapters.NewGoogleTokenVerifier(googleClientIDs)
	useCase := application.NewGoogleLoginUseCase(repo, verifier, jwtSecret)
	return controllers.NewGoogleLoginController(useCase, cookieSecure)
}
