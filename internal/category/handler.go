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

	categoryGroup := router.Group("/categories")
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
