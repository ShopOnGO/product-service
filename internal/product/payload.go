package product

import (
	"github.com/ShopOnGO/product-service/internal/productVariant"
)

type BaseProductEvent struct {
	Action string 				 `json:"action"`
	Product ProductCreatedEvent  `json:"product"`
}

type ProductCreatedEvent struct {
	Name        string 	`json:"name"`
	Description string 	`json:"description"`
	Material    string 	`gorm:"type:varchar(200)"`
	IsActive    bool   	`json:"is_active"`

	CategoryID  uint   	`json:"category_id"`
	BrandID     uint   	`json:"brand_id"`

	ImageKeys  []string `json:"image_keys"`
	VideoKeys  []string `json:"video_keys"`

	Variants []productVariant.ProductVariant `json:"variants"`
}

type ProductCreatedEventForMediaAndSearch struct {
	Action    string 	`json:"action"`
	ProductID uint   	`json:"product_id"`

	// Полные данные продукта — для Search Service
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Material    string 	`gorm:"type:varchar(200)"`
	IsActive    bool    `json:"is_active"`
	CategoryID  uint    `json:"category_id"`
	BrandID     uint    `json:"brand_id"`

	// Данные для Media Service
	ImageKeys []string 	`json:"image_keys"`
	VideoKeys []string 	`json:"video_keys"`

	Variants []productVariant.ProductVariant `json:"variants"`
}

// type ProductCreatedEventForMedia struct {
// 	Action    string 	`json:"action"`
// 	ProductID uint   	`json:"product_id"`
// 	ImageKeys []string 	`json:"image_keys"`
// 	VideoKeys []string 	`json:"video_keys"`
// }

type MediaUpdateEvent struct {
	Action    string 	`json:"action"`
	ProductID uint   	`json:"product_id"`
    ImageURLs []string 	`json:"image_urls"`
	VideoURLs []string 	`json:"video_urls"`
}
