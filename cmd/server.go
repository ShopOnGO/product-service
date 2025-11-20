package main

import (
	"context"
	"fmt"
	"net"

	GoogleGRPC "google.golang.org/grpc"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	pb "github.com/ShopOnGO/product-proto/pkg/product"
	"github.com/ShopOnGO/product-service/configs"
	"github.com/ShopOnGO/product-service/internal/brand"
	"github.com/ShopOnGO/product-service/internal/category"
	"github.com/ShopOnGO/product-service/internal/grpc"
	"github.com/ShopOnGO/product-service/internal/product"
	"github.com/ShopOnGO/product-service/internal/productVariant"
	"github.com/ShopOnGO/product-service/migrations"
	"github.com/ShopOnGO/product-service/pkg/db"
	"github.com/segmentio/kafka-go"

	"github.com/gin-gonic/gin"

	// Пустой импорт _ говорит Go-компилятору не удалять его,
	// а для swag это является прямой командой проанализировать пакет.
	_ "github.com/ShopOnGO/review-proto/pkg/service"
)

// @title Product Service API
// @version 1.0
// @description API для управления продуктами, категориями и брендами.
// @host localhost:8082
// @BasePath /
// @schemes http
func main() {
	migrations.CheckForMigrations()
	conf := configs.LoadConfig()
	consoleLvl := conf.LogLevel
	fileLvl := conf.FileLogLevel
	logger.InitLogger(consoleLvl, fileLvl)
	logger.EnableFileLogging("TailorNado_product-service")

	database := db.NewDB(conf)
	kafkaProducers := kafkaService.InitKafkaProducers(
		conf.KafkaProducer.Brokers,
		conf.KafkaProducer.Topic,
	)
	router := gin.Default()

	// repository
	productRepo := product.NewProductRepository(database)
	brandRepo := brand.NewBrandRepository(database)
	categoryRepo := category.NewCategoryRepository(database)
	productVariantRepo := productVariant.NewProductVariantRepository(database)

	// service
	productService := product.NewProductService(productRepo)
	brandService := brand.NewBrandService(brandRepo)
	categoryService := category.NewCategoryService(categoryRepo)
	productVariantService := productVariant.NewProductVariantService(productVariantRepo)

	// handler
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductSvc: productService,
		Kafka:      kafkaProducers["products"],
	})
	brand.NewBrandHandler(router, brand.BrandHandlerDeps{
		BrandSvc: brandService,
		Kafka:    kafkaProducers["brands"],
	})
	category.NewCategoryHandler(router, category.CategoryHandlerDeps{
		CategorySvc: categoryService,
		Kafka:       kafkaProducers["categories"],
	})
	productVariant.NewProductVariantHandler(router, productVariant.ProductVariantHandlerDeps{
		ProductVariantSvc: productVariantService,
		Kafka:             kafkaProducers["variants"],
	})
	grpc.NewReviewHandler(router)

	kafkaProductConsumer := kafkaService.NewConsumer(
		conf.KafkaProduct.Brokers,
		conf.KafkaProduct.Topic,
		conf.KafkaProduct.GroupID,
		conf.KafkaProduct.ClientID,
	)
	kafkaVariantConsumer := kafkaService.NewConsumer(
		conf.KafkaVariant.Brokers,
		conf.KafkaVariant.Topic,
		conf.KafkaVariant.GroupID,
		conf.KafkaVariant.ClientID,
	)
	kafkaMediaConsumer := kafkaService.NewConsumer(
		conf.KafkaMedia.Brokers,
		conf.KafkaMedia.Topic,
		conf.KafkaMedia.GroupID,
		conf.KafkaMedia.ClientID,
	)

	defer kafkaProductConsumer.Close()
	defer kafkaVariantConsumer.Close()
	defer kafkaMediaConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go kafkaProductConsumer.Consume(ctx, func(msg kafka.Message) error {
		key := string(msg.Key)
		return product.HandleProductEvent(msg.Value, key, productService, productVariantService, kafkaProducers["products"])
	})

	go kafkaVariantConsumer.Consume(ctx, func(msg kafka.Message) error {
		key := string(msg.Key)
		return productVariant.HandleProductVariantEvent(msg.Value, key, productVariantService)
	})
	go kafkaMediaConsumer.Consume(ctx, func(msg kafka.Message) error {
		key := string(msg.Key)
		return product.HandleProductEvent(msg.Value, key, productService, productVariantService, nil)
	})

	go func() {
		listener, err := net.Listen("tcp", ":50053")
		if err != nil {
			logger.Infof("TCP listener error: %v\n", err)
			return
		}

		grpcServer := GoogleGRPC.NewServer()
		pb.RegisterProductVariantServiceServer(grpcServer, productVariant.NewGrpcProductVariantService(productVariantService))
		pb.RegisterProductServiceServer(grpcServer, product.NewGrpcProductService(productService))

		logger.Info("gRPC server listening on :50053")
		if err := grpcServer.Serve(listener); err != nil {
			logger.Infof("gRPC server error: %v\n", err)
		}
	}()

	go func() {
		if err := router.Run(":8082"); err != nil {
			fmt.Println("Ошибка при запуске HTTP-сервера:", err)
		}
	}()

	select {}
}
