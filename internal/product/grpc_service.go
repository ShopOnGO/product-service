package product

import (
	"context"

	pb "github.com/ShopOnGO/product-proto/pkg/product"
)

type GrpcProductService struct {
	pb.UnimplementedProductServiceServer
	productSvc *ProductService
}

func NewGrpcProductService(svc *ProductService) *GrpcProductService {
	return &GrpcProductService{productSvc: svc}
}

func (g *GrpcProductService) GetProductsByIDs(ctx context.Context, req *pb.GetProductsByIDsRequest) (*pb.GetProductsByIDsResponse, error) {
	productIDs := make([]uint, len(req.ProductIds))
	for i, id := range req.ProductIds {
		productIDs[i] = uint(id)
	}

	products, err := g.productSvc.GetProductsByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	grpcProducts := make([]*pb.Product, len(products))
	for i, p := range products {
		grpcProducts[i] = &pb.Product{
			Id:           uint64(p.ID),
			Name:         p.Name,
			Description:  p.Description,
			Rating:       p.Rating.InexactFloat64(),
			ReviewCount:  uint32(p.ReviewCount),
			RatingSum:    uint32(p.RatingSum),
			QuestionCount:uint32(p.QuestionCount),
			IsActive:     p.IsActive,
			CategoryId:   uint32(p.CategoryID),
			BrandId:      uint32(p.BrandID),
			ImageUrls:    p.ImageURLs,
			VideoUrls:    p.VideoURLs,
		}
	}

	return &pb.GetProductsByIDsResponse{
		Products: grpcProducts,
	}, nil
}

