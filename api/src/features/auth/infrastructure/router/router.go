package router

import (
	"net/http"

	"vault/src/features/auth/infrastructure/controllers"
)

func RegisterRoutes(mux *http.ServeMux, login *controllers.LoginController, googleLogin *controllers.GoogleLoginController) {
	mux.HandleFunc("POST /api/v1/auth/login", login.Handle)
	mux.HandleFunc("POST /api/v1/auth/google", googleLogin.Handle)
}
