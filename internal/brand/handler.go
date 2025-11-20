package brand

import (
	"net/http"
	"strconv"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BrandHandler struct {
	brandSvc *BrandService
	Kafka    *kafkaService.KafkaService // Добавлено
}

// Новая структура для зависимостей
type BrandHandlerDeps struct {
	BrandSvc *BrandService
	Kafka    *kafkaService.KafkaService
}

func NewBrandHandler(router *gin.Engine, deps BrandHandlerDeps) *BrandHandler {
	handler := &BrandHandler{
		brandSvc: deps.BrandSvc,
		Kafka:    deps.Kafka,
	}

	brandGroup := router.Group("/product-service/brands")
	{
		brandGroup.GET("/", handler.GetBrands)
		brandGroup.GET("/:id", handler.GetBrandByID)
		brandGroup.POST("/", handler.CreateBrand)
		brandGroup.PUT("/:id", handler.UpdateBrand)
		brandGroup.DELETE("/:id", handler.DeleteBrand)
	}

	return handler
}

// GetBrands godoc
// @Summary Получить список всех брендов
// @Description Возвращает все бренды
// @Tags Бренды
// @Success 200 {array} brand.Brand
// @Failure 500 {object} gin.H "Ошибка при получении брендов"
// @Router /brands/ [get]
func (h *BrandHandler) GetBrands(c *gin.Context) {
	brands, err := h.brandSvc.GetAllBrands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch brands"})
		return
	}
	c.JSON(http.StatusOK, brands)
}

// GetBrandByID godoc
// @Summary Получить бренд по ID
// @Description Возвращает бренд по его уникальному идентификатору
// @Tags Бренды
// @Param id path int true "ID бренда"
// @Success 200 {object} brand.Brand
// @Failure 400 {object} gin.H "Некорректный ID бренда"
// @Failure 404 {object} gin.H "Бренд не найден"
// @Router /brands/{id} [get]
func (h *BrandHandler) GetBrandByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}
	brand, err := h.brandSvc.GetBrandByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, brand)
}

// CreateBrand godoc
// @Summary Создать новый бренд
// @Description Создаёт новый бренд на основе переданных данных
// @Tags Бренды
// @Accept json
// @Produce json
// @Param brand body brand.BrandRequest true "Данные для создания бренда"
// @Success 201 {object} brand.Brand
// @Failure 400 {object} gin.H "Некорректный формат запроса"
// @Failure 500 {object} gin.H "Ошибка при создании бренда"
// @Router /brands/ [post]
func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var payload BrandRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	newBrand := &Brand{
		Name:        payload.Name,
		Description: payload.Description,
		VideoURL:    payload.VideoURL,
		Logo:        payload.Logo,
	}

	createdBrand, err := h.brandSvc.CreateBrand(newBrand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create brand"})
		return
	}

	// go h.sendNotification(
	// 	c,
	// 	"notification-BrandCreated", // Kafka Key
	// 	"BRAND_CREATED",             // Category
	// 	"brand",                     // Subtype
	// 	map[string]interface{}{ // Payload
	// 		"brandID":   createdBrand.ID,
	// 		"brandName": createdBrand.Name,
	// 		"message":   fmt.Sprintf("Новый бренд '%s' был успешно создан.", createdBrand.Name),
	// 	},
	// )
	c.JSON(http.StatusCreated, createdBrand)
}

// UpdateBrand godoc
// @Summary Обновить бренд
// @Description Обновляет существующий бренд по ID
// @Tags Бренды
// @Accept json
// @Produce json
// @Param id path int true "ID бренда"
// @Param brand body brand.BrandRequest true "Данные для обновления бренда"
// @Success 200 {object} brand.Brand
// @Failure 400 {object} gin.H "Некорректный формат запроса"
// @Failure 500 {object} gin.H "Ошибка при обновлении бренда"
// @Router /brands/{id} [put]
func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}

	var payload BrandRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	newBrand := &Brand{
		Model:       gorm.Model{ID: uint(id)},
		Name:        payload.Name,
		Description: payload.Description,
		VideoURL:    payload.VideoURL,
		Logo:        payload.Logo,
	}

	updatedBrand, err := h.brandSvc.UpdateBrand(newBrand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// go h.sendNotification(
	// 	c,
	// 	"notification-BrandUpdated", // Kafka Key
	// 	"BRAND_UPDATED",             // Category
	// 	"brand",                     // Subtype
	// 	map[string]interface{}{ // Payload
	// 		"brandID":   updatedBrand.ID,
	// 		"brandName": updatedBrand.Name,
	// 		"message":   fmt.Sprintf("Бренд '%s' был обновлен.", updatedBrand.Name),
	// 	},
	// )
	c.JSON(http.StatusOK, updatedBrand)
}

// DeleteBrand godoc
// @Summary Удалить бренд
// @Description Удаляет бренд по ID
// @Tags Бренды
// @Param id path int true "ID бренда"
// @Success 200 {object} gin.H "Сообщение об успешном удалении"
// @Failure 400 {object} gin.H "Некорректный ID бренда"
// @Failure 500 {object} gin.H "Ошибка при удалении бренда"
// @Router /brands/{id} [delete]
func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}

	if err := h.brandSvc.DeleteBrand(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// go h.sendNotification(
	// 	c,
	// 	"notification-BrandDeleted", // Kafka Key
	// 	"BRAND_DELETED",             // Category
	// 	"brand",                     // Subtype
	// 	map[string]interface{}{ // Payload
	// 		"brandID": id,
	// 		"message": fmt.Sprintf("Бренд '%v' был удалён.", id),
	// 	},
	// )
	c.JSON(http.StatusOK, gin.H{"message": "brand deleted"})
}

// func (h *BrandHandler) sendNotification(
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
// 		return // Не прерываем основной запрос из-за ошибки уведомления
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
