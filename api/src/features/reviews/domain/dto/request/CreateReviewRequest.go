package request

import "errors"

type CreateReviewRequest struct {
	ProviderID string `json:"provider_id"`
	Content    string `json:"content"`
}

func (r CreateReviewRequest) Validate() error {
	if r.ProviderID == "" {
		return errors.New("el proveedor es obligatorio")
	}
	if r.Content == "" {
		return errors.New("el contenido es obligatorio")
	}
	return nil
}
