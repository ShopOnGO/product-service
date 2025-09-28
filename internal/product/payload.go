package product

import (
	"github.com/ShopOnGO/product-service/internal/productVariant"
	"github.com/shopspring/decimal"
)

type BaseProductEvent struct {
	Action string 				 `json:"action"`
	Product ProductCreatedEvent  `json:"product"`
}

type ProductCreatedEvent struct {
	Name        	string 			`json:"name"`
	Description 	string 			`json:"description"`
	Material    	string 			`gorm:"type:varchar(200)"`
	IsActive    	bool   			`json:"is_active"`
	Rating      	decimal.Decimal `gorm:"type:decimal(8,1);not null;default:0"`
	ReviewCount   	uint      		`gorm:"not null;default:0"`
	RatingSum     	uint	  		`gorm:"not null;default:0"`
	QuestionCount	uint 			`gorm:"default:0"`

	CategoryID  	uint   			`json:"category_id"`
	BrandID     	uint   			`json:"brand_id"`

	ImageKeys  		[]string 		`json:"image_keys"`
	VideoKeys  		[]string 		`json:"video_keys"`

	Variants 		[]productVariant.ProductVariant `json:"variants"`
}


type ProductCreatedEventForMediaAndSearch struct {
	Action    	string 		`json:"action"`
	ProductID 	uint   		`json:"product_id"`

	// Полные данные продукта — для Search Service
	Name        string 		`json:"name"`
	Description string  	`json:"description"`
	Material    string 		`gorm:"type:varchar(200)"`
	Rating      float64 	`gorm:"type:decimal(8,1);not null;default:0"`
	ReviewCount	uint   		`gorm:"not null;default:0"`
	IsActive    bool    	`json:"is_active"`
	CategoryID  uint    	`json:"category_id"`
	BrandID     uint    	`json:"brand_id"`

	// Данные для Media Service
	ImageKeys 	[]string 	`json:"image_keys"`
	VideoKeys 	[]string 	`json:"video_keys"`

	Variants 	[]*ProductVariantForEvent `json:"variants"`
}

type ProductVariantForEvent struct {
    VariantID      uint     `json:"variant_id"`
    SKU            string   `json:"sku"`
    Price          float64  `json:"price"`
    Discount       float64  `json:"discount"`
    Sizes          string   `json:"sizes"`
    Colors         string   `json:"colors"`
    Stock          uint32   `json:"stock"`
    Barcode        string   `json:"barcode,omitempty"`
    Dimensions     string   `json:"dimensions,omitempty"`
    ImageURLs      []string `json:"image_urls,omitempty"`
    MinOrder       uint     `json:"min_order,omitempty"`
    IsActive       bool     `json:"is_active"`
    ReservedStock  uint32   `json:"reserved_stock,omitempty"`
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
