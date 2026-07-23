package entities

import "time"

type ChatMessage struct {
	ID              string
	SenderID        string
	RecipientID     string
	CipherText      string
	EncryptedAESKey string
	IV              string
	Status          string
	CreatedAt       time.Time
}
