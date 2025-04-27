package brand

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BrandHandler struct {
	brandSvc *BrandService
}

func NewBrandHandler(router *gin.Engine, brandSvc *BrandService) *BrandHandler {
	handler := &BrandHandler{brandSvc: brandSvc}
	
	brandGroup := router.Group("/brands")
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
	c.JSON(http.StatusOK, gin.H{"message": "brand deleted"})
}