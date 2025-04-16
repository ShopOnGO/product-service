package brand

type BrandRequest struct {
	ID          uint   `json:"id,omitempty"` // для update или delete
	Name        string `json:"name"`
	Description string `json:"description"`
	VideoURL    string `json:"video_url"`
	Logo        string `json:"logo"`
}