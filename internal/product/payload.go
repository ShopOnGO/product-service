package product

type BaseProductEvent struct {
	Action string 				 `json:"action"`
	Product ProductCreatedEvent  `json:"product"`
}

type ProductCreatedEvent struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Discount    int64  `json:"discount"`
	IsActive    bool   `json:"is_active"`

	CategoryID  uint   `json:"category_id"`
	BrandID     uint   `json:"brand_id"`

	ImageKeys  []string `json:"image_keys"`
	VideoKeys  []string `json:"video_keys"`
}

type ProductCreatedEventForMedia struct {
	Action    	string   `json:"action"`
	ProductID 	uint     `json:"product_id"`
	ImageKeys  	[]string `json:"image_keys"`
	VideoKeys  	[]string `json:"video_keys"`
}

type MediaUpdateEvent struct {
	Action    string 	`json:"action"`
	ProductID uint   	`json:"product_id"`
    ImageURLs []string 	`json:"image_urls"`
	VideoURLs []string 	`json:"video_urls"`
}
