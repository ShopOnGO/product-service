package migrations

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
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
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true, //временно игнорировать миграции в первый раз а потом их добавить
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&product.Product{}, 
		&productVariant.ProductVariant{}, 
		&category.Category{}, 
		&brand.Brand{})
		
	if err != nil {
		return err
	}

	logger.Info("✅")
	return nil
}
