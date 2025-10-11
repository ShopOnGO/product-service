package migrations

import (
	"fmt"
	"os"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/ShopOnGO/product-service/configs"
	"github.com/ShopOnGO/product-service/internal/brand"
	"github.com/ShopOnGO/product-service/internal/category"
	"github.com/ShopOnGO/product-service/internal/product"
	"github.com/ShopOnGO/product-service/internal/productVariant"
)

func CheckForMigrations() error {

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		logger.Info("🚀 Starting migrations...")
		if err := RunMigrations(); err != nil {
			logger.Errorf("Error processing migrations: %v", err)
		}
		return nil
	}
	// if not "migrate" args[1]
	return nil
}

func RunMigrations() error {
	// Попробуем загрузить .env, но не паникуем, если его нет
	cfg := configs.LoadConfig()

	dsn := os.Getenv("DSN")
	if dsn == "" {
		return fmt.Errorf("DSN is empty, check your .env or environment variables")
	}

	db, err := gorm.Open(postgres.Open(cfg.Db.Dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	// Авто-миграции
	if err := db.AutoMigrate(
		&product.Product{},
		&productVariant.ProductVariant{},
		&category.Category{},
		&brand.Brand{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("✅ Migrations completed")
	return nil
}
