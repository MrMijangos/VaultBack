package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/chat/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	sendChatMessage *controllers.SendChatMessageController,
	getConversationMessages *controllers.GetConversationMessagesController,
	updateChatMessageStatus *controllers.UpdateChatMessageStatusController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/chat/messages", auth(http.HandlerFunc(sendChatMessage.Handle)))
	mux.Handle("GET /api/v1/conversations/{id}/messages", auth(http.HandlerFunc(getConversationMessages.Handle)))
	mux.Handle("PATCH /api/v1/chat/messages/{id}/status", auth(http.HandlerFunc(updateChatMessageStatus.Handle)))
}
