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

func (h *BrandHandler) GetBrands(c *gin.Context) {
	brands, err := h.brandSvc.GetAllBrands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch brands"})
		return
	}
	c.JSON(http.StatusOK, brands)
}

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