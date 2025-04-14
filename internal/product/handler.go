package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productSvc *ProductService
}

func NewProductHandler(router *gin.Engine, productSvc *ProductService) *ProductHandler {
	handler := &ProductHandler{productSvc: productSvc}

	productGroup := router.Group("/products")
	{
		productGroup.GET("/", handler.GetProducts)
		productGroup.GET("/:id", handler.GetProductByID)
		productGroup.POST("/", handler.CreateProduct)
		productGroup.PUT("/:id", handler.UpdateProduct)
		productGroup.DELETE("/:id", handler.DeleteProduct)
	}

	return handler
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productSvc.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.productSvc.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	product, err := h.productSvc.CreateProduct(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var updated Product
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	product, err := h.productSvc.UpdateProduct(uint(id), &updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	if err := h.productSvc.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

func HandleProductEvent(msg []byte, key string, productSvc *ProductService) error {

	logger.Infof("Получено сообщение: %s", string(msg))
	
	var base BaseProductEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	switch base.Action {
	case "create":
		event := base.Product
		logger.Infof("Получены данные для создания продукта: name=%q, category_id=%d, brand_id=%d, price=%d",
			event.Name, event.CategoryID, event.BrandID, event.Price)

		logger.Infof("Получены данные для создания продукта: %+v", event)
		newProduct := &Product{
			Name:        event.Name,
			Description: event.Description,
			Price:       event.Price,
			Discount:    event.Discount,
			IsActive:    event.IsActive,
			CategoryID:  event.CategoryID,
			BrandID:     event.BrandID,
			Images:      event.Images,
			VideoURL:    event.VideoURL,
		}

		createdProduct, err := productSvc.CreateProduct(newProduct)
		if err != nil {
			logger.Errorf("Ошибка при создании отзыва: %v", err)
			return err
		}
		logger.Infof("Продукт успешно создан: %+v", createdProduct)
	default:
		return fmt.Errorf("неизвестное действие для продукта: %s", base.Action)
	}
	return nil
}
