package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	categorySvc *CategoryService
}

func NewCategoryHandler(router *gin.Engine, categorySvc *CategoryService) *CategoryHandler {
	handler := &CategoryHandler{categorySvc: categorySvc}

	categoryGroup := router.Group("/product-service/categories")
	{
		categoryGroup.POST("/", handler.CreateCategory)
		categoryGroup.GET("/featured", handler.GetFeaturedCategories)
		categoryGroup.GET("/by-name", handler.GetCategoryByName)
		categoryGroup.GET("/:id", handler.GetCategoryByID)
		categoryGroup.PUT("/:id", handler.UpdateCategory)
		categoryGroup.DELETE("/:id", handler.DeleteCategory)
	}

	return handler
}

// CreateCategory godoc
// @Summary Создать новую категорию
// @Description Создает новую категорию с заданными данными
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CategoryPayload true "Данные категории"
// @Success 201 {object} Category
// @Failure 400 {object} gin.H "Неверный запрос"
// @Failure 500 {object} gin.H "Ошибка сервера"
// @Router /categories/ [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var payload CategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	category := &Category{
		Name:             payload.Name,
		Description:      payload.Description,
		ImageURL:         payload.ImageURL,
		ParentCategoryID: payload.ParentCategoryID,
	}

	created, err := h.categorySvc.CreateCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetFeaturedCategories godoc
// @Summary Получить популярные категории
// @Description Получает несколько популярных категорий
// @Tags categories
// @Accept json
// @Produce json
// @Param amount query int false "Количество категорий (по умолчанию 10)"
// @Success 200 {array} Category
// @Failure 500 {object} gin.H "Ошибка сервера"
// @Router /categories/featured [get]
func (h *CategoryHandler) GetFeaturedCategories(c *gin.Context) {
	amountStr := c.Query("amount")
	amount := 10 // default
	if amountStr != "" {
		if parsed, err := strconv.Atoi(amountStr); err == nil {
			amount = parsed
		}
	}

	categories, err := h.categorySvc.GetFeaturedCategories(amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID godoc
// @Summary Получить категорию по ID
// @Description Ищет категорию по её идентификатору
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} Category
// @Failure 400 {object} gin.H "Неверный ID"
// @Failure 404 {object} gin.H "Категория не найдена"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	category, err := h.categorySvc.GetCategoryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryByName godoc
// @Summary Получить категорию по названию
// @Description Ищет категорию по её имени
// @Tags categories
// @Accept json
// @Produce json
// @Param name query string true "Название категории"
// @Success 200 {object} Category
// @Failure 400 {object} gin.H "Отсутствует параметр name"
// @Failure 404 {object} gin.H "Категория не найдена"
// @Router /categories/by-name [get]
func (h *CategoryHandler) GetCategoryByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'name' is required"})
		return
	}

	category, err := h.categorySvc.GetCategoryByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}


// UpdateCategory godoc
// @Summary Обновить категорию
// @Description Обновляет существующую категорию по её ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "ID категории"
// @Param category body CategoryPayload true "Обновленные данные категории"
// @Success 200 {object} Category
// @Failure 400 {object} gin.H "Неверный запрос"
// @Failure 500 {object} gin.H "Ошибка сервера"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var payload CategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	category := &Category{
		Model:            gorm.Model{ID: uint(id)},
		Name:             payload.Name,
		Description:      payload.Description,
		ImageURL:         payload.ImageURL,
		ParentCategoryID: payload.ParentCategoryID,
	}

	updated, err := h.categorySvc.UpdateCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteCategory godoc
// @Summary Удалить категорию
// @Description Удаляет категорию по её ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} gin.H "Категория удалена"
// @Failure 400 {object} gin.H "Неверный ID"
// @Failure 500 {object} gin.H "Ошибка сервера"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.categorySvc.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}
