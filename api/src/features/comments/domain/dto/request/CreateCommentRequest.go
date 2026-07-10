package request

import "errors"

type CreateCommentRequest struct {
	Content string `json:"content"`
}

func (r CreateCommentRequest) Validate() error {
	if r.Content == "" {
		return errors.New("el contenido es obligatorio")
	}
	return nil
}
