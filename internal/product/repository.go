package product

import (
	"errors"

	"github.com/ShopOnGO/product-service/pkg/db"
	"gorm.io/gorm"
)

// ProductRepository предоставляет методы для работы с продуктами в базе данных.
type ProductRepository struct {
	Db *db.Db
}

func NewProductRepository(db *db.Db) *ProductRepository {
	return &ProductRepository{
		Db: db,
	}
}

func (r *ProductRepository) GetAll() ([]Product, error) {
	var products []Product
	if err := r.Db.
		Preload("Category").
		Preload("Brand").
		Preload("Variants").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id uint) (*Product, error) {
	var product Product
	if err := r.Db.
		Preload("Category").
		Preload("Brand").
		Preload("Variants").
		First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}


func (r *ProductRepository) Create(product *Product) error {
	if err := r.Db.Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) Update(product *Product) error {
	if err := r.Db.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) Delete(id uint) error {
	if err := r.Db.Delete(&Product{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) IsProductOwnedByUser(productID, userID uint) (bool, error) {
    var prod Product
    err := r.Db.
        Select("id").
        Where("id = ? AND seller_id = ?", productID, userID).
        First(&prod).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }

    return true, nil
}