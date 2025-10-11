package grpc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type ReviewHandler struct {
	Clients *GRPCClients
}

func InitGRPCClients() *GRPCClients {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("REVIEW_SERVICE_HOST"), os.Getenv("REVIEW_SERVICE_PORT")), grpc.WithInsecure())
	if err != nil {
		logger.Errorf("Ошибка подключения к gRPC серверу: %v", err)
	}

	logger.Info("gRPC connected")
	reviewClient := pb.NewReviewServiceClient(conn)
	questionClient := pb.NewQuestionServiceClient(conn)

	return &GRPCClients{
		ReviewClient:   reviewClient,
		QuestionClient: questionClient,
	}
}

func NewReviewHandler(router *gin.Engine) {
	handler := &ReviewHandler{
		Clients: InitGRPCClients(),
	}

	productGroup := router.Group("/products")
	{
		productGroup.GET("/reviews/:id", handler.GetProductWithReviews)
		productGroup.GET("/questions/:id", handler.GetProductWithQuestions)
	}
}

// GetProductWithReviews получает отзывы о продукте
// @Summary Получение отзывов по ID варианта продукта
// @Description Возвращает список отзывов для заданного варианта продукта
// @Tags Отзывы
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Param limit query int false "Количество отзывов для получения"
// @Param offset query int false "Смещение для пагинации"
// @Success 200 {object} pb.ReviewListResponse
// @Failure 400 {object} map[string]string "Неверный ID продукта"
// @Failure 500 {object} map[string]string "Ошибка получения отзывов"
// @Router /products/reviews/{id} [get]
func (h *ReviewHandler) GetProductWithReviews(c *gin.Context) {
	ctx := context.Background()
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	productVariantID, err := strconv.Atoi(c.Param("id"))
	if err != nil || productVariantID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_variant_id"})
		return
	}

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	resp, err := h.Clients.ReviewClient.GetReviewsForProduct(ctx, &pb.GetReviewsRequest{
		ProductVariantId: uint32(productVariantID),
		Limit:            int32(limit),
		Offset:           int32(offset),
	})
	if err != nil {
		logger.Errorf("Ошибка при вызове gRPC GetReviewsForProduct: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить отзывы"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProductWithQuestions получает вопросы о продукте
// @Summary Получение вопросов по ID варианта продукта
// @Description Возвращает список вопросов для заданного варианта продукта
// @Tags Вопросы
// @Accept json
// @Produce json
// @Param id path int true "ID варианта продукта"
// @Param limit query int false "Количество вопросов для получения"
// @Param offset query int false "Смещение для пагинации"
// @Success 200 {object} pb.QuestionListResponse
// @Failure 400 {object} map[string]string "Неверный ID продукта"
// @Failure 500 {object} map[string]string "Ошибка получения вопросов"
// @Router /products/questions/{id} [get]
func (h *ReviewHandler) GetProductWithQuestions(c *gin.Context) {
	ctx := context.Background()
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	productVariantID, err := strconv.Atoi(c.Param("id"))
	if err != nil || productVariantID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_variant_id"})
		return
	}

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	resp, err := h.Clients.QuestionClient.GetQuestionsForProduct(ctx, &pb.GetQuestionsRequest{
		ProductVariantId: uint32(productVariantID),
		Limit:            int32(limit),
		Offset:           int32(offset),
	})
	if err != nil {
		logger.Errorf("Ошибка при вызове gRPC GetQuestionsForProduct: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить вопросы"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
