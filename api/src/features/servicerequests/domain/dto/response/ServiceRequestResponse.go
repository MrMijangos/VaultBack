package response

import (
	"time"

	"vault/src/features/servicerequests/domain/entities"
)

type ServiceRequestResponse struct {
	ID            string     `json:"id"`
	AssetID       string     `json:"asset_id"`
	AssetName     string     `json:"asset_name"`
	AssetImageURL string     `json:"asset_image_url"`
	OwnerID       string     `json:"owner_id"`
	OwnerName     string     `json:"owner_name"`
	BusinessID    string     `json:"business_id"`
	BusinessName  string     `json:"business_name"`
	Type          string     `json:"type"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	AcceptedAt    *time.Time `json:"accepted_at"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	ConfirmedAt   *time.Time `json:"confirmed_at"`
}

func FromEntity(sr entities.ServiceRequest) ServiceRequestResponse {
	return ServiceRequestResponse{
		ID:            sr.ID,
		AssetID:       sr.AssetID,
		AssetName:     sr.AssetName,
		AssetImageURL: sr.AssetImageURL,
		OwnerID:       sr.OwnerID,
		OwnerName:     sr.OwnerName,
		BusinessID:    sr.BusinessID,
		BusinessName:  sr.BusinessName,
		Type:          sr.Type,
		Status:        sr.Status,
		CreatedAt:     sr.CreatedAt,
		AcceptedAt:    sr.AcceptedAt,
		StartedAt:     sr.StartedAt,
		FinishedAt:    sr.FinishedAt,
		ConfirmedAt:   sr.ConfirmedAt,
	}
}

func FromEntities(list []entities.ServiceRequest) []ServiceRequestResponse {
	out := make([]ServiceRequestResponse, 0, len(list))
	for _, sr := range list {
		out = append(out, FromEntity(sr))
	}
	return out
}
