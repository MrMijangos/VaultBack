package response

import (
	"time"

	"vault/src/features/addresses/domain/entities"
)

type AddressResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Label      string    `json:"label"`
	Recipient  string    `json:"recipient"`
	Phone      string    `json:"phone"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	References string    `json:"references"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
}

func FromEntity(a entities.Address) AddressResponse {
	return AddressResponse{
		ID:         a.ID,
		UserID:     a.UserID,
		Label:      a.Label,
		Recipient:  a.Recipient,
		Phone:      a.Phone,
		Street:     a.Street,
		City:       a.City,
		State:      a.State,
		PostalCode: a.PostalCode,
		References: a.ReferenceNotes,
		IsDefault:  a.IsDefault,
		CreatedAt:  a.CreatedAt,
	}
}

func FromEntities(list []entities.Address) []AddressResponse {
	out := make([]AddressResponse, 0, len(list))
	for _, a := range list {
		out = append(out, FromEntity(a))
	}
	return out
}
