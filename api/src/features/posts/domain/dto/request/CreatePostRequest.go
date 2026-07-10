package request

import "errors"

type CreatePostRequest struct {
	Content string `json:"content"`
	AssetID string `json:"asset_id"`
}

func (r CreatePostRequest) Validate() error {
	if r.Content == "" {
		return errors.New("el contenido es obligatorio")
	}
	return nil
}
