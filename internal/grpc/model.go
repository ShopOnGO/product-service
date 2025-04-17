package grpc

import (
	pb "github.com/ShopOnGO/review-proto/pkg/service"
)

type GRPCClients struct {
	ReviewClient	pb.ReviewServiceClient
	QuestionClient	pb.QuestionServiceClient
}