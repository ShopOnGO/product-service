package productVariant

import "github.com/shopspring/decimal"

type CreateProductVariantPayload struct {
	ProductID     uint     			`json:"product_id" binding:"required"`
	SKU           string   			`json:"sku" binding:"required"`
	Price    	  decimal.Decimal   `json:"price" binding:"required"`
	Discount 	  decimal.Decimal   `json:"discount"`
	ReservedStock uint32   			`json:"reserved_stock"`
	Sizes  		  []uint32 			`json:"sizes" binding:"omitempty"`
	Colors        []string 			`json:"colors" binding:"omitempty"`
	Stock         uint32   			`json:"stock"`
	Material      string   			`json:"material"`
	Barcode       string   			`json:"barcode"`
	IsActive      bool     			`json:"is_active"`
	Images        []string 			`json:"images" binding:"omitempty"`
	MinOrder      uint     			`json:"min_order"`
	Dimensions    string   			`json:"dimensions"`
}

type ProductVariantCreatedEvent struct {
	SKU           string   			`json:"sku" binding:"required"`
	Price    	  decimal.Decimal   `json:"price" binding:"required"`
	Discount 	  decimal.Decimal   `json:"discount"`
	ReservedStock uint32   			`json:"reserved_stock"`
	Sizes  		  []uint32 			`json:"sizes" binding:"omitempty"`
	Colors        []string 			`json:"colors" binding:"omitempty"`
	Stock         uint32   			`json:"stock"`
	Material      string   			`json:"material"`
	Barcode       string   			`json:"barcode"`
	IsActive      bool     			`json:"is_active"`
	Images        []string 			`json:"images" binding:"omitempty"`
	MinOrder      uint     			`json:"min_order"`
	Dimensions    string   			`json:"dimensions"`
}

type BaseProductVariantEvent struct {
	Action  		string                   	`json:"action"`
	ProductID		uint					 	`json:"product_id"`
	ProductVariant 	ProductVariantCreatedEvent 	`json:"product_variant"`
	UserID 			uint                     	`json:"user_id"`
}

type UpdateProductVariantPayload struct {
	Price    	  *decimal.Decimal 	`json:"price" gorm:"type:decimal(8,2);not null"`
	Discount 	  *decimal.Decimal 	`json:"discount" gorm:"type:decimal(8,2);not null;default:0"`
	ReservedStock *uint32          	`json:"reserved_stock"`
	Sizes         *[]uint32        	`json:"sizes"`
	Colors        *[]string        	`json:"colors"`
	Stock         *uint32          	`json:"stock"`
	Material      *string          	`json:"material"`
	Barcode       *string          	`json:"barcode"`
	IsActive      *bool            	`json:"is_active"`
	Images        *[]string        	`json:"images"`
	MinOrder      *uint            	`json:"min_order"`
	Dimensions    *string          	`json:"dimensions"`
}

type ReserveStockPayload struct {
	Quantity uint32 `json:"quantity" binding:"required,gt=0"`
}

type ReleaseStockPayload struct {
	Quantity uint32 `json:"quantity" binding:"required,gt=0"`
}

type UpdateStockPayload struct {
	Stock uint32 `json:"stock" binding:"required"`
}