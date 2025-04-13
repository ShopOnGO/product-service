package main

import (
	"context"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/product-service/configs"
	"github.com/ShopOnGO/product-service/internal/product"
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

	// Инициализация Kafka-консьюмера
	kafkaConsumer := kafkaService.NewConsumer(
		conf.Kafka.Brokers,
		conf.Kafka.Topic,
		conf.Kafka.GroupID,
		conf.Kafka.ClientID,
	)
	defer kafkaConsumer.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go kafkaConsumer.Consume(ctx, func(msg kafka.Message) error {
		return gin.Error{}
	})

	// repository
	productRepo := product.NewProductRepository(database)

	// service
	productService := product.NewProductService(productRepo)

	// handler
	product.NewProductHandler(router, productService)

	go func() {
		if err := router.Run(":8082"); err != nil {
			fmt.Println("Ошибка при запуске HTTP-сервера:", err)
		}
	}()
	
	select{}
}

