package main

import (
	"context"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/product-service/configs"
	"github.com/ShopOnGO/product-service/internal/brand"
	"github.com/ShopOnGO/product-service/internal/category"
	"github.com/ShopOnGO/product-service/internal/grpc"
	"github.com/ShopOnGO/product-service/internal/product"
	"github.com/ShopOnGO/product-service/internal/productVariant"
	"github.com/ShopOnGO/product-service/migrations"
	"github.com/ShopOnGO/product-service/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	migrations.CheckForMigrations()
	conf := configs.LoadConfig()
	database := db.NewDB(conf)
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
	product.NewProductHandler(router, productService)
	brand.NewBrandHandler(router, brandService)
	category.NewCategoryHandler(router, categoryService)
	productVariant.NewProductVariantHandler(router, productVariantService)
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

	defer kafkaProductConsumer.Close()
	defer kafkaVariantConsumer.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go kafkaProductConsumer.Consume(ctx, func(msg kafka.Message) error {
		key := string(msg.Key)
		return product.HandleProductEvent(msg.Value, key, productService)
	})
	go kafkaVariantConsumer.Consume(ctx, func(msg kafka.Message) error {
		key := string(msg.Key)
		return productVariant.HandleProductVariantEvent(msg.Value, key, productVariantService)
	})


	go func() {
		if err := router.Run(":8082"); err != nil {
			fmt.Println("Ошибка при запуске HTTP-сервера:", err)
		}
	}()
	
	select{}
}

