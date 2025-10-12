package product

import (
	"errors"

	"gorm.io/gorm"
)

type ProductService struct {
	repo *ProductRepository
}

func NewProductService(repository *ProductRepository) *ProductService {
	return &ProductService{
		repo: repository,
	}
}

func (s *ProductService) GetAllProducts() ([]Product, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductByID(id uint) (*Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetProductsByIDs(ids []uint) ([]Product, error) {
	products, err := s.repo.GetProductsByIDs(ids)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) CreateProduct(product *Product) (*Product, error) {
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(id uint, updated *Product) (*Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Обновляем только нужные поля
	product.Name = updated.Name
	product.Description = updated.Description
	product.Material = updated.Material
	product.IsActive = updated.IsActive
	product.CategoryID = updated.CategoryID
	product.BrandID = updated.BrandID
	product.ImageURLs = updated.ImageURLs
	product.VideoURLs = updated.VideoURLs

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) UpdateProductMedia(productID uint, images []string, video []string) error {
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return err
	}

	product.ImageURLs = images
	product.VideoURLs = video

	return s.repo.Update(product)
}


func (s *ProductService) DeleteProduct(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
