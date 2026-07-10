package request

import "errors"

type UpdatePostRequest struct {
	Content string `json:"content"`
}

func (r UpdatePostRequest) Validate() error {
	if r.Content == "" {
		return errors.New("el contenido es obligatorio")
	}
	return nil
}
