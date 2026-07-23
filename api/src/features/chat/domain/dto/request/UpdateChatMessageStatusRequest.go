package request

import "errors"

type UpdateChatMessageStatusRequest struct {
	Status string `json:"status"`
}

func (r UpdateChatMessageStatusRequest) Validate() error {
	if r.Status != "delivered" && r.Status != "read" {
		return errors.New("el estado debe ser 'delivered' o 'read'")
	}
	return nil
}
