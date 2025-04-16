package productVariant

import (
	"net/http"
	"strconv"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ProductVariantHandler содержит зависимость от сервиса вариантов продукта.
type ProductVariantHandler struct {
	productVariantSvc productVariantService
}

// NewProductVariantHandler регистрирует маршруты для работы с вариантами продукта.
func NewProductVariantHandler(router *gin.Engine, productVariantSvc productVariantService) *ProductVariantHandler {
	handler := &ProductVariantHandler{productVariantSvc: productVariantSvc}

	variantGroup := router.Group("/product-variants")
	{
		variantGroup.GET("/:id", handler.GetProductVariantByID)
		variantGroup.GET("/by-sku", handler.GetProductVariantBySKU)
		variantGroup.POST("/", handler.CreateProductVariant)
		variantGroup.PUT("/:id", handler.UpdateProductVariant)
		variantGroup.DELETE("/:id", handler.DeleteProductVariant)
		
		variantGroup.POST("/:id/reserve", handler.ReserveStock)
		variantGroup.POST("/:id/release", handler.ReleaseStock)
		variantGroup.PUT("/:id/stock", handler.UpdateStock)
		variantGroup.GET("/:id/available", handler.GetAvailableStock)	
	}

	return handler
}

// CreateProductVariant создает новый вариант продукта.
func (h *ProductVariantHandler) CreateProductVariant(c *gin.Context) {
	var payload CreateProductVariantPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	variant := &ProductVariant{
		ProductID:     payload.ProductID,
		SKU:           payload.SKU,
		Price:         payload.Price,
		Discount:      payload.Discount,
		ReservedStock: payload.ReservedStock,
		Rating:        payload.Rating,
		Sizes:         payload.Sizes,
		Colors:        payload.Colors,
		Stock:         payload.Stock,
		Material:      payload.Material,
		Barcode:       payload.Barcode,
		IsActive:      payload.IsActive,
		Images:        payload.Images,
		MinOrder:      payload.MinOrder,
		Dimensions:    payload.Dimensions,
	}

	created, err := h.productVariantSvc.CreateProductVariant(variant)
	if err != nil {
		logger.Errorf("Error creating product variant: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetProductVariantByID возвращает вариант продукта по его ID.
func (h *ProductVariantHandler) GetProductVariantByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	variant, err := h.productVariantSvc.GetProductVariantByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, variant)
}

// GetProductVariantBySKU возвращает вариант продукта по артикулу.
// Запрос: GET /variants/by-sku?sku=SKU123
func (h *ProductVariantHandler) GetProductVariantBySKU(c *gin.Context) {
	sku := c.Query("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'sku' is required"})
		return
	}

	variant, err := h.productVariantSvc.GetBySKU(sku)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if variant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product variant not found"})
		return
	}
	c.JSON(http.StatusOK, variant)
}

func (h *ProductVariantHandler) UpdateProductVariant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	var payload UpdateProductVariantPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	updated, err := h.productVariantSvc.UpdateProductVariantByInput(uint(id), payload)
	if err != nil {
		logger.Errorf("Error updating product variant: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteProductVariant выполняет мягкое удаление варианта продукта.
func (h *ProductVariantHandler) DeleteProductVariant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	if err := h.productVariantSvc.DeleteProductVariant(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product variant deleted"})
}


// stock
func (h *ProductVariantHandler) ReserveStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	var payload ReserveStockPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
		return
	}

	if err := h.productVariantSvc.ReserveStock(uint(id), payload.Quantity); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock reserved"})
}

func (h *ProductVariantHandler) ReleaseStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	var payload ReleaseStockPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
		return
	}

	if err := h.productVariantSvc.ReleaseStock(uint(id), payload.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock released"})
}

func (h *ProductVariantHandler) UpdateStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	var payload UpdateStockPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock value"})
		return
	}

	if err := h.productVariantSvc.UpdateStock(uint(id), payload.Stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock updated"})
}

func (h *ProductVariantHandler) GetAvailableStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	available, err := h.productVariantSvc.GetAvailableStock(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available_stock": available})
}