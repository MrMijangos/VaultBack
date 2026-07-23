package request

import "errors"

// SendChatMessageRequest lleva un paquete E2EE opaco -- ninguno de estos
// campos se interpreta como texto en el servidor, solo se valida que no
// esten vacios.
type SendChatMessageRequest struct {
	RecipientID     string `json:"recipient_id"`
	CipherText      string `json:"cipher_text"`
	EncryptedAESKey string `json:"encrypted_aes_key"`
	IV              string `json:"iv"`
}

func (r SendChatMessageRequest) Validate() error {
	if r.RecipientID == "" {
		return errors.New("el destinatario es obligatorio")
	}
	if r.CipherText == "" {
		return errors.New("el texto cifrado es obligatorio")
	}
	if r.EncryptedAESKey == "" {
		return errors.New("la llave cifrada es obligatoria")
	}
	if r.IV == "" {
		return errors.New("el vector de inicializacion es obligatorio")
	}
	return nil
}
