package app

import (
    "context"
    "net"

    "github.com/gin-gonic/gin"
    "github.com/segmentio/kafka-go"
    googleGrpc "google.golang.org/grpc"

    "github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
    "github.com/ShopOnGO/ShopOnGO/pkg/logger"
    "github.com/ShopOnGO/product-service/configs"
    "github.com/ShopOnGO/product-service/internal/brand"
    "github.com/ShopOnGO/product-service/internal/category"
    "github.com/ShopOnGO/product-service/internal/product"
    "github.com/ShopOnGO/product-service/internal/productVariant"
    "github.com/ShopOnGO/product-service/internal/grpc"
    "github.com/ShopOnGO/product-service/migrations"
    "github.com/ShopOnGO/product-service/pkg/db"
    pb "github.com/ShopOnGO/product-proto/pkg/product"
)

type App struct {
    Router               *gin.Engine
    ProductService       *product.ProductService
    BrandService         *brand.BrandService
    CategoryService      *category.CategoryService
    VariantService       *productVariant.ProductVariantService
    KafkaProductConsumer *kafkaService.KafkaService
    KafkaVariantConsumer *kafkaService.KafkaService
    GrpcServer           *googleGrpc.Server
}

func NewApp() *App {
    migrations.CheckForMigrations()
    conf := configs.LoadConfig()
    database := db.NewDB(conf)

    // Репозитории
    prodRepo := product.NewProductRepository(database)
    brandRepo := brand.NewBrandRepository(database)
    categoryRepo := category.NewCategoryRepository(database)
    variantRepo := productVariant.NewProductVariantRepository(database)

    // Сервисы
    prodSvc := product.NewProductService(prodRepo)
    brandSvc := brand.NewBrandService(brandRepo)
    categorySvc := category.NewCategoryService(categoryRepo)
    variantSvc := productVariant.NewProductVariantService(variantRepo)

    // роутер и хендлеры
    router := gin.Default()
    product.NewProductHandler(router, prodSvc)
    brand.NewBrandHandler(router, brandSvc)
    category.NewCategoryHandler(router, categorySvc)
    productVariant.NewProductVariantHandler(router, variantSvc)
	grpc.NewReviewHandler(router)

    // Kafka-консьюмеры
    kafkaProd := kafkaService.NewConsumer(
        conf.KafkaProduct.Brokers,
        conf.KafkaProduct.Topic,
        conf.KafkaProduct.GroupID,
        conf.KafkaProduct.ClientID,
    )
    kafkaVar := kafkaService.NewConsumer(
        conf.KafkaVariant.Brokers,
        conf.KafkaVariant.Topic,
        conf.KafkaVariant.GroupID,
        conf.KafkaVariant.ClientID,
    )

    // gRPC-сервер
    grpcServer := googleGrpc.NewServer()
    pb.RegisterProductVariantServiceServer(
        grpcServer,
        productVariant.NewGrpcProductVariantService(variantSvc),
    )

    return &App{
        Router:               router,
        ProductService:       prodSvc,
        BrandService:         brandSvc,
        CategoryService:      categorySvc,
        VariantService:       variantSvc,
        KafkaProductConsumer: kafkaProd,
        KafkaVariantConsumer: kafkaVar,
        GrpcServer:           grpcServer,
    }
}

func (a *App) RunKafka(ctx context.Context) {
    defer a.KafkaProductConsumer.Close()
    defer a.KafkaVariantConsumer.Close()

    go a.KafkaProductConsumer.Consume(ctx, func(msg kafka.Message) error {
        return product.HandleProductEvent(msg.Value, string(msg.Key), a.ProductService)
    })
    go a.KafkaVariantConsumer.Consume(ctx, func(msg kafka.Message) error {
        return productVariant.HandleProductVariantEvent(msg.Value, string(msg.Key), a.VariantService)
    })
}

func (a *App) RunGRPC() {
    listener, err := net.Listen("tcp", ":50053")
    if err != nil {
        logger.Errorf("TCP listener error: %v", err)
        return
    }
    logger.Info("gRPC server listening on :50053")
    if err := a.GrpcServer.Serve(listener); err != nil {
        logger.Errorf("gRPC server error: %v", err)
    }
}

func (a *App) RunHTTP() {
    addr := ":8082"
    logger.Infof("HTTP server listening on %s", addr)
    if err := a.Router.Run(addr); err != nil {
        logger.Errorf("HTTP server error: %v", err)
    }
}
