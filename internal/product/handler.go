package product

import (
	"net/http"
	"strconv"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/gin-gonic/gin"
)

type ProductHandlerDeps struct {
	ProductSvc  *ProductService
	Kafka 		*kafkaService.KafkaService
}

type ProductHandler struct {
	ProductSvc  *ProductService
	Kafka 		*kafkaService.KafkaService
}

func NewProductHandler(router *gin.Engine, deps ProductHandlerDeps) *ProductHandler {
	handler := &ProductHandler{
		ProductSvc:    	deps.ProductSvc,
		Kafka: 			deps.Kafka,
	}

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

// GetProducts получает все продукты
// @Summary Получение всех продуктов
// @Description Возвращает список всех продуктов
// @Tags Продукты
// @Accept json
// @Produce json
// @Success 200 {array} Product
// @Failure 500 {object} map[string]string "Ошибка получения продуктов"
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.ProductSvc.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductByID получает продукт по ID
// @Summary Получение продукта по ID
// @Description Возвращает продукт по его ID
// @Tags Продукты
// @Accept json
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} Product
// @Failure 400 {object} map[string]string "Неверный ID продукта"
// @Failure 404 {object} map[string]string "Продукт не найден"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.ProductSvc.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct создаёт новый продукт
// @Summary Создание нового продукта
// @Description Создаёт новый продукт с переданными данными
// @Tags Продукты
// @Accept json
// @Produce json
// @Param product body Product true "Данные продукта"
// @Success 201 {object} Product
// @Failure 400 {object} map[string]string "Неверное тело запроса"
// @Failure 500 {object} map[string]string "Ошибка при создании продукта"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	product, err := h.ProductSvc.CreateProduct(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct обновляет данные продукта
// @Summary Обновление данных продукта
// @Description Обновляет информацию о продукте по его ID
// @Tags Продукты
// @Accept json
// @Produce json
// @Param id path int true "ID продукта"
// @Param product body Product true "Данные для обновления продукта"
// @Success 200 {object} Product
// @Failure 400 {object} map[string]string "Неверное тело запроса"
// @Failure 404 {object} map[string]string "Продукт не найден"
// @Failure 500 {object} map[string]string "Ошибка при обновлении продукта"
// @Router /products/{id} [put]
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

	product, err := h.ProductSvc.UpdateProduct(uint(id), &updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct удаляет продукт
// @Summary Удаление продукта
// @Description Удаляет продукт по его ID
// @Tags Продукты
// @Accept json
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} map[string]string "Продукт удалён"
// @Failure 400 {object} map[string]string "Неверный ID продукта"
// @Failure 500 {object} map[string]string "Ошибка при удалении продукта"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	if err := h.ProductSvc.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

