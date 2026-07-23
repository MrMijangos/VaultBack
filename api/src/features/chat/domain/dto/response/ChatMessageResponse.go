package response

import (
	"time"

	"vault/src/features/chat/domain/entities"
)

type ChatMessageResponse struct {
	ID              string    `json:"id"`
	SenderID        string    `json:"sender_id"`
	RecipientID     string    `json:"recipient_id"`
	CipherText      string    `json:"cipher_text"`
	EncryptedAESKey string    `json:"encrypted_aes_key"`
	IV              string    `json:"iv"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

func FromEntity(m entities.ChatMessage) ChatMessageResponse {
	return ChatMessageResponse{
		ID:              m.ID,
		SenderID:        m.SenderID,
		RecipientID:     m.RecipientID,
		CipherText:      m.CipherText,
		EncryptedAESKey: m.EncryptedAESKey,
		IV:              m.IV,
		Status:          m.Status,
		CreatedAt:       m.CreatedAt,
	}
}

func FromEntities(list []entities.ChatMessage) []ChatMessageResponse {
	out := make([]ChatMessageResponse, 0, len(list))
	for _, m := range list {
		out = append(out, FromEntity(m))
	}
	return out
}
