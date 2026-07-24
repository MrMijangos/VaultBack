package entities

// ConversationSummary es una fila de la bandeja de chat: el otro
// participante, su último mensaje (todavía cifrado -- el servidor nunca ve
// el texto plano) y cuántos mensajes suyos siguen sin marcarse como
// leídos.
type ConversationSummary struct {
	OtherUserID        string
	OtherUserName      string
	OtherUserAvatarURL string
	LastMessage        ChatMessage
	UnreadCount        int
}
