package category

import (
	"errors"

	"github.com/ShopOnGO/product-service/pkg/db"
)

type CategoryRepository struct {
	Db *db.Db
}

func NewCategoryRepository(db *db.Db) *CategoryRepository {
	return &CategoryRepository{
		Db: db,
	}
}

func (repo *CategoryRepository) Create(category *Category) (*Category, error) {
	result := repo.Db.Create(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (repo *CategoryRepository) GetFeaturedCategories(amount int) ([]Category, error) {
	if amount > 20 {
		amount = 20
	}
	var categories []Category
	query := repo.Db.DB.Preload("SubCategories")

	if amount > 0 {
		query = query.Limit(amount)
	}

	result := query.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (repo *CategoryRepository) GetByName(name string) (*Category, error) {
	var category Category
	result := repo.Db.Preload("SubCategories").First(&category, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func (repo *CategoryRepository) GetByID(id uint) (*Category, error) {
	if id == 0 {
		return nil, errors.New("invalid category ID")
	}
	var category Category
	result := repo.Db.Preload("SubCategories").First(&category, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func (repo *CategoryRepository) Update(category *Category) (*Category, error) {
	result := repo.Db.Model(&Category{}).Where("id = ?", category.ID).Updates(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (repo *CategoryRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid  ID")
	}
	result := repo.Db.Delete(&Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

