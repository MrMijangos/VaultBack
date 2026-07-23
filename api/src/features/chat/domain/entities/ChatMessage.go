package entities

import "time"

type ChatMessage struct {
	ID              string
	SenderID        string
	RecipientID     string
	CipherText      string
	EncryptedAESKey string
	// EncryptedAESKeySender es la misma llave AES del mensaje, cifrada
	// ademas con la publica del propio emisor -- sin esto, quien envia un
	// mensaje no podria releerlo despues (su privada no destraba una llave
	// cifrada para la publica del receptor).
	EncryptedAESKeySender string
	IV                    string
	Status                string
	CreatedAt             time.Time
}
