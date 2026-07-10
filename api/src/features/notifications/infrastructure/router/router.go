package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/notifications/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createNotification *controllers.CreateNotificationController,
	getMyNotifications *controllers.GetMyNotificationsController,
	markAsRead *controllers.MarkNotificationAsReadController,
	deleteNotification *controllers.DeleteNotificationController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/notifications", auth(http.HandlerFunc(createNotification.Handle)))
	mux.Handle("GET /api/v1/notifications", auth(http.HandlerFunc(getMyNotifications.Handle)))
	mux.Handle("PUT /api/v1/notifications/{id}/read", auth(http.HandlerFunc(markAsRead.Handle)))
	mux.Handle("DELETE /api/v1/notifications/{id}", auth(http.HandlerFunc(deleteNotification.Handle)))
}
