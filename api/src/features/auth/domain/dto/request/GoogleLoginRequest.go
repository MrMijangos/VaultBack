package request

import "errors"

type GoogleLoginRequest struct {
	IDToken string `json:"id_token"`
}

func (r GoogleLoginRequest) Validate() error {
	if r.IDToken == "" {
		return errors.New("falta el token de Google")
	}
	return nil
}
