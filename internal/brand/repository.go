package brand

import (
	"errors"

	"github.com/ShopOnGO/product-service/pkg/db"
)

type BrandRepository struct {
	Db *db.Db
}

func NewBrandRepository(db *db.Db) *BrandRepository {
	return &BrandRepository{
		Db: db,
	}
}

func (repo *BrandRepository) GetByID(id uint) (*Brand, error) {
	var brand Brand
	if err := repo.Db.First(&brand, id).Error; err != nil {
		return nil, err
	}
	return &brand, nil
}

func (repo *BrandRepository) GetAll() ([]*Brand, error) {
	var brands []*Brand
	if err := repo.Db.Find(&brands).Error; err != nil {
		return nil, err
	}
	return brands, nil
}

func (repo *BrandRepository) Create(brand *Brand) (*Brand, error) {
	result := repo.Db.Create(brand)
	if result.Error != nil {
		return nil, result.Error
	}
	return brand, nil
}

func (repo *BrandRepository) Update(brand *Brand) (*Brand, error) {
	result := repo.Db.Model(&Brand{}).Where("id = ?", brand.ID).Updates(brand)
	if result.Error != nil {
		return nil, result.Error
	}
	return brand, nil
}

func (repo *BrandRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid brand ID")
	}
	result := repo.Db.Delete(&Brand{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

