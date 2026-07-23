package request

import "errors"

type CreateAssetCommentRequest struct {
	Content string `json:"content"`
}

func (r CreateAssetCommentRequest) Validate() error {
	if r.Content == "" {
		return errors.New("el contenido es obligatorio")
	}
	return nil
}
