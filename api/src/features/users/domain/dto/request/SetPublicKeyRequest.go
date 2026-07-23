package request

import "errors"

type SetPublicKeyRequest struct {
	PublicKey string `json:"public_key"`
}

func (r SetPublicKeyRequest) Validate() error {
	if r.PublicKey == "" {
		return errors.New("la llave publica es obligatoria")
	}
	return nil
}
