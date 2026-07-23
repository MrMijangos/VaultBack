package response

type PublicKeyResponse struct {
	UserID    string  `json:"user_id"`
	PublicKey *string `json:"public_key"`
}
