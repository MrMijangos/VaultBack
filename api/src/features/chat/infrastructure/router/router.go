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
	getConversations *controllers.GetConversationsController,
	deleteChatMessage *controllers.DeleteChatMessageController,
	deleteConversation *controllers.DeleteConversationController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/chat/messages", auth(http.HandlerFunc(sendChatMessage.Handle)))
	mux.Handle("GET /api/v1/conversations", auth(http.HandlerFunc(getConversations.Handle)))
	mux.Handle("GET /api/v1/conversations/{id}/messages", auth(http.HandlerFunc(getConversationMessages.Handle)))
	mux.Handle("DELETE /api/v1/conversations/{id}", auth(http.HandlerFunc(deleteConversation.Handle)))
	mux.Handle("PATCH /api/v1/chat/messages/{id}/status", auth(http.HandlerFunc(updateChatMessageStatus.Handle)))
	mux.Handle("DELETE /api/v1/chat/messages/{id}", auth(http.HandlerFunc(deleteChatMessage.Handle)))
}
