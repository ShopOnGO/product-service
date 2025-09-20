package product

import (
	"github.com/ShopOnGO/product-service/internal/brand"
	"github.com/ShopOnGO/product-service/internal/category"
	"github.com/ShopOnGO/product-service/internal/productVariant"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        string 				`gorm:"type:varchar(255);not null" json:"name"`
	Description string 				`gorm:"type:text" json:"description"`
	Material    string 				`gorm:"type:varchar(200)"`
	IsActive    bool   				`gorm:"default:true" json:"is_active"`

	CategoryID 	uint              	`gorm:"not null" json:"category_id"`
	Category   	category.Category 	`gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`

	BrandID 	uint        		`gorm:"not null" json:"brand_id"`
	Brand   	brand.Brand 		`gorm:"foreignKey:BrandID;constraint:OnDelete:CASCADE"`

	Variants 	[]productVariant.ProductVariant `gorm:"foreignKey:ProductID"`

	ImageURLs 	pq.StringArray 		`gorm:"type:text[]"`
    VideoURLs 	pq.StringArray 		`gorm:"type:text[]"`
}