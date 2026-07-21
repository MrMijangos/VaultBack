package request

type UpdateAdRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	TargetSection string `json:"target_section"`
	TargetID      string `json:"target_id"`
}
