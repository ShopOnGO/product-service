package productVariant

import (
	"net/http"
	"strconv"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Новая структура для зависимостей
type ProductVariantHandlerDeps struct {
	ProductVariantSvc *ProductVariantService
	Kafka             *kafkaService.KafkaService
}

type ProductVariantHandler struct {
	productVariantSvc *ProductVariantService
	Kafka             *kafkaService.KafkaService // Добавлено
}

func NewProductVariantHandler(router *gin.Engine, deps ProductVariantHandlerDeps) *ProductVariantHandler {
	handler := &ProductVariantHandler{
		productVariantSvc: deps.ProductVariantSvc,
		Kafka:             deps.Kafka,
	}

	variantGroup := router.Group("/product-service/product-variants")
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
// @Summary Создание нового варианта продукта
// @Description Создает новый вариант продукта с указанными данными.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param variant body CreateProductVariantPayload true "Данные для создания варианта продукта"
// @Success 201 {object} ProductVariant
// @Failure 400 {object} map[string]string "Неверное тело запроса"
// @Failure 500 {object} map[string]string "Ошибка при создании варианта продукта"
// @Router /product-variants [post]
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
		Sizes:         payload.Sizes,
		Colors:        payload.Colors,
		Stock:         payload.Stock,
		Barcode:       payload.Barcode,
		IsActive:      payload.IsActive,
		ImageURLs:     payload.Images,
		MinOrder:      payload.MinOrder,
		Dimensions:    payload.Dimensions,
	}

	created, err := h.productVariantSvc.CreateProductVariant(variant)
	if err != nil {
		logger.Errorf("Error creating product variant: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// go h.sendNotification(
	// 	c,
	// 	"notification-VariantCreated", // Kafka Key
	// 	"VARIANT_CREATED",             // Category
	// 	"product_variant",             // Subtype
	// 	map[string]interface{}{ // Payload
	// 		"variantID":  created.ID,
	// 		"productID":  created.ProductID,
	// 		"variantSKU": created.SKU,
	// 		"message":    fmt.Sprintf("Новый вариант '%s' для товара %d был успешно создан.", created.SKU, created.ProductID),
	// 	},
	// )
	c.JSON(http.StatusCreated, created)
}

// GetProductVariantByID получает вариант продукта по ID.
// @Summary Получение варианта продукта по ID
// @Description Возвращает информацию о варианте продукта по его ID.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Success 200 {object} ProductVariant
// @Failure 400 {object} map[string]string "Неверный ID варианта продукта"
// @Failure 404 {object} map[string]string "Вариант продукта не найден"
// @Router /product-variants/{id} [get]
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

// GetProductVariantBySKU получает вариант продукта по артикулу.
// @Summary Получение варианта продукта по артикулу
// @Description Возвращает вариант продукта по артикулу SKU.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param sku query string true "SKU варианта продукта"
// @Success 200 {object} ProductVariant
// @Failure 400 {object} map[string]string "Не указан SKU"
// @Failure 404 {object} map[string]string "Вариант продукта не найден"
// @Router /product-variants/by-sku [get]
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

// UpdateProductVariant обновляет данные варианта продукта.
// @Summary Обновление данных варианта продукта
// @Description Обновляет информацию о варианте продукта по ID.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Param variant body UpdateProductVariantPayload true "Данные для обновления варианта продукта"
// @Success 200 {object} ProductVariant
// @Failure 400 {object} map[string]string "Неверное тело запроса"
// @Failure 404 {object} map[string]string "Вариант продукта не найден"
// @Failure 500 {object} map[string]string "Ошибка при обновлении варианта продукта"
// @Router /product-variants/{id} [put]
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

	// go h.sendNotification(
	// 	c,
	// 	"notification-VariantUpdated", // Kafka Key
	// 	"VARIANT_UPDATED",             // Category
	// 	"product_variant",             // Subtype
	// 	map[string]interface{}{ // Payload
	// 		"variantID":  updated.ID,
	// 		"productID":  updated.ProductID,
	// 		"variantSKU": updated.SKU,
	// 		"message":    fmt.Sprintf("Вариант '%s' (товар %d) был обновлен.", updated.SKU, updated.ProductID),
	// 	},
	// )
	c.JSON(http.StatusOK, updated)
}

// DeleteProductVariant выполняет мягкое удаление варианта продукта.
// @Summary Удаление варианта продукта
// @Description Удаляет вариант продукта по ID.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Success 200 {object} map[string]string "Вариант продукта удален"
// @Failure 400 {object} map[string]string "Неверный ID варианта продукта"
// @Failure 500 {object} map[string]string "Ошибка при удалении варианта продукта"
// @Router /product-variants/{id} [delete]
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

// ReserveStock резервирует товар на складе.
// @Summary Резервирование товара
// @Description Резервирует указанное количество товара для варианта продукта.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Param quantity body ReserveStockPayload true "Количество для резервирования"
// @Success 200 {object} map[string]string "Запас зарезервирован"
// @Failure 400 {object} map[string]string "Неверное количество"
// @Failure 409 {object} map[string]string "Ошибка при резервировании товара"
// @Router /product-variants/{id}/reserve [post]
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

// ReleaseStock освобождает зарезервированный товар.
// @Summary Освобождение товара
// @Description Освобождает зарезервированное количество товара для варианта продукта.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Param quantity body ReleaseStockPayload true "Количество для освобождения"
// @Success 200 {object} map[string]string "Запас освобожден"
// @Failure 400 {object} map[string]string "Неверное количество"
// @Failure 500 {object} map[string]string "Ошибка при освобождении товара"
// @Router /product-variants/{id}/release [post]
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

// UpdateStock обновляет запас товара.
// @Summary Обновление запаса товара
// @Description Обновляет количество товара на складе для варианта продукта.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Param stock body UpdateStockPayload true "Новое количество товара"
// @Success 200 {object} map[string]string "Запас обновлен"
// @Failure 400 {object} map[string]string "Неверное количество товара"
// @Failure 500 {object} map[string]string "Ошибка при обновлении запаса"
// @Router /product-variants/{id}/stock [put]
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

// GetAvailableStock получает доступный запас товара.
// @Summary Получение доступного запаса товара
// @Description Возвращает доступный запас для варианта продукта.
// @Tags Варианты Продуктов
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Success 200 {object} map[string]int "Доступный запас"
// @Failure 400 {object} map[string]string "Неверный ID варианта продукта"
// @Failure 500 {object} map[string]string "Ошибка при получении доступного запаса"
// @Router /product-variants/{id}/available [get]
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

// func (h *ProductVariantHandler) sendNotification(
// 	c *gin.Context,
// 	kafkaKey string,
// 	category string,
// 	subtype string,
// 	payload map[string]interface{},
// ) {
// 	// 1. Получаем userID из контекста
// 	rawUserID, exists := c.Get("userID")
// 	if !exists {
// 		log.Printf("⚠️ [Kafka] userID не найден в контексте для %s, уведомление не отправлено", category)
// 		return
// 	}

// 	userID, ok := rawUserID.(uint32)
// 	if !ok {
// 		log.Printf("⚠️ [Kafka] userID в контексте имеет неверный тип для %s, уведомление не отправлено", category)
// 		return
// 	}

// 	// 2. Создаем тело уведомления (JSON-контракт)
// 	notificationPayload := map[string]interface{}{
// 		"category": category,
// 		"subtype":  subtype,
// 		"userID":   userID,
// 		"payload":  payload,
// 	}

// 	// 3. Маршалим в JSON
// 	jsonPayload, err := json.Marshal(notificationPayload)
// 	if err != nil {
// 		log.Printf("🚨 [Kafka] Ошибка маршалинга уведомления %s: %v", category, err)
// 		return
// 	}

// 	// 4. Публикуем сообщение
// 	if err := h.Kafka.Produce(c, []byte(kafkaKey), jsonPayload); err != nil {
// 		log.Printf("🚨 [Kafka] Не удалось опубликовать сообщение %s: %v", category, err)
// 	} else {
// 		log.Printf("✅ [Kafka] Уведомление %s отправлено для userID %d", category, userID)
// 	}
// }
