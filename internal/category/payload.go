package category

type CategoryPayload struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	ImageURL         string `json:"image_url"`
	ParentCategoryID *uint  `json:"parent_category_id"`
}