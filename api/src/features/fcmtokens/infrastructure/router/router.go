package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/fcmtokens/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	registerFCMToken *controllers.RegisterFCMTokenController,
	deleteFCMToken *controllers.DeleteFCMTokenController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/users/fcm-token", auth(http.HandlerFunc(registerFCMToken.Handle)))
	mux.Handle("DELETE /api/v1/users/fcm-token", auth(http.HandlerFunc(deleteFCMToken.Handle)))
}
