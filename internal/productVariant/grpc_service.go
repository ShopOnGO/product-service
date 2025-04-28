package productVariant

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	pb "github.com/ShopOnGO/product-proto/pkg/product"
	"gorm.io/gorm"
)

type GrpcProductVariantService struct {
    pb.UnimplementedProductVariantServiceServer
    productVariantSvc *ProductVariantService
}

func NewGrpcProductVariantService(svc *ProductVariantService) *GrpcProductVariantService {
    return &GrpcProductVariantService{productVariantSvc: svc}
}

func (g *GrpcProductVariantService) CheckProductVariantExists(ctx context.Context, req *pb.CheckProductVariantRequest) (*pb.CheckProductVariantResponse, error) {
    variant, err := g.productVariantSvc.GetProductVariantByID(uint(req.ProductVariantId))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &pb.CheckProductVariantResponse{Exists: false, IsActive: false}, nil
        }
        logger.Errorf("CheckProductVariantExists error: %v", err)
        return nil, fmt.Errorf("internal error checking variant: %w", err)
    }

    return &pb.CheckProductVariantResponse{
        Exists:   true,
        IsActive: variant.IsActive,
    }, nil
}

func (g *GrpcProductVariantService) GetProductVariants(ctx context.Context, req *pb.GetProductVariantsRequest) (*pb.GetProductVariantsResponse, error) {
    ids := make([]uint, len(req.ProductVariantIds))
    for i, id := range req.ProductVariantIds {
        ids[i] = uint(id)
    }

    variants, err := g.productVariantSvc.GetVariantsByIDs(ids)
    if err != nil {
        logger.Errorf("GetProductVariants error: %v", err)
        return nil, fmt.Errorf("internal error fetching variants: %w", err)
    }

    resp := &pb.GetProductVariantsResponse{}
    for _, v := range variants {
        resp.ProductVariants = append(resp.ProductVariants, &pb.ProductVariant{
            Id:            uint64(v.ID),
            ProductId:     uint64(v.ProductID),
            Sku:           v.SKU,
            Price:         v.Price.String(),
            Discount:      v.Discount.String(),
            IsActive:      v.IsActive,
            Stock:         uint32(v.Stock),
            Images:        v.Images,
            Rating:        float64(v.Rating.InexactFloat64()),
            ReviewCount:   uint32(v.ReviewCount),
        })
    }

    return resp, nil
}
