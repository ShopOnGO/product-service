package grpc

import (
	"context"
	"net/http"
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
	conn, err := grpc.Dial("review_container:50052", grpc.WithInsecure())
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

