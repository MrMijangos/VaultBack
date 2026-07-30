package request

import "errors"

type RegisterFCMTokenRequest struct {
	Token    string  `json:"token"`
	Platform *string `json:"platform"`
}

func (r RegisterFCMTokenRequest) Validate() error {
	if r.Token == "" {
		return errors.New("el token es obligatorio")
	}
	return nil
}

type DeleteFCMTokenRequest struct {
	Token string `json:"token"`
}

func (r DeleteFCMTokenRequest) Validate() error {
	if r.Token == "" {
		return errors.New("el token es obligatorio")
	}
	return nil
}
