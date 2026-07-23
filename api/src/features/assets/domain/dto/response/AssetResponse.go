package response

import (
	"time"

	"vault/src/features/assets/domain/entities"
)

type AssetPhotoResponse struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	IsCover bool   `json:"is_cover"`
	Order   int    `json:"order"`
}

type AssetResponse struct {
	ID              string               `json:"id"`
	UserID          string               `json:"user_id"`
	Name            string               `json:"name"`
	Category        string               `json:"category"`
	Brand           string               `json:"brand"`
	PurchaseValue   *float64             `json:"purchase_value"`
	Condition       string               `json:"condition"`
	PurchaseDate    *string              `json:"purchase_date"`
	StoreOrigin     string               `json:"store_origin"`
	Notes           string               `json:"notes"`
	BlockchainTxID  string               `json:"blockchain_tx_id"`
	BlockchainHash  string               `json:"blockchain_hash"`
	CreatedAt       time.Time            `json:"created_at"`
	Photos          []AssetPhotoResponse `json:"photos"`
	IsForSale       bool                 `json:"is_for_sale"`
	SalePrice       *float64             `json:"sale_price"`
	SaleDescription string               `json:"sale_description"`
	Size            string               `json:"size"`
	SellerName      string               `json:"seller_name"`
	SellerAvatarURL string               `json:"seller_avatar_url"`
}

func FromEntity(asset entities.Asset, photos []entities.AssetPhoto) AssetResponse {
	var purchaseDate *string
	if asset.PurchaseDate != nil {
		formatted := asset.PurchaseDate.Format("2006-01-02")
		purchaseDate = &formatted
	}

	photoResponses := make([]AssetPhotoResponse, 0, len(photos))
	for _, p := range photos {
		photoResponses = append(photoResponses, AssetPhotoResponse{
			ID:      p.ID,
			URL:     p.URL,
			IsCover: p.IsCover,
			Order:   p.Order,
		})
	}

	return AssetResponse{
		ID:              asset.ID,
		UserID:          asset.UserID,
		Name:            asset.Name,
		Category:        asset.Category,
		Brand:           asset.Brand,
		PurchaseValue:   asset.PurchaseValue,
		Condition:       asset.Condition,
		PurchaseDate:    purchaseDate,
		StoreOrigin:     asset.StoreOrigin,
		Notes:           asset.Notes,
		BlockchainTxID:  asset.BlockchainTxID,
		BlockchainHash:  asset.BlockchainHash,
		CreatedAt:       asset.CreatedAt,
		Photos:          photoResponses,
		IsForSale:       asset.IsForSale,
		SalePrice:       asset.SalePrice,
		SaleDescription: asset.SaleDescription,
		Size:            asset.Size,
		SellerName:      asset.SellerName,
		SellerAvatarURL: asset.SellerAvatarURL,
	}
}

func FromEntities(assets []entities.Asset) []AssetResponse {
	list := make([]AssetResponse, 0, len(assets))
	for _, a := range assets {
		list = append(list, FromEntity(a, nil))
	}
	return list
}
