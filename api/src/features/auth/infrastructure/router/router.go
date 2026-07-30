package router

import (
	"net/http"

	"vault/src/core/middleware"
	"vault/src/features/auth/infrastructure/controllers"
)

func RegisterRoutes(mux *http.ServeMux, login *controllers.LoginController, googleLogin *controllers.GoogleLoginController) {
	mux.Handle("POST /api/v1/auth/login", middleware.RateLimitLogin(http.HandlerFunc(login.Handle)))
	mux.HandleFunc("POST /api/v1/auth/google", googleLogin.Handle)
}
