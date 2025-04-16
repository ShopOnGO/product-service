package brand

import "errors"

type BrandService struct {
	repo *BrandRepository
}

func NewBrandService(repository *BrandRepository) *BrandService {
	return &BrandService{
		repo: repository,
	}
}

func (s *BrandService) GetBrandByID(id uint) (*Brand, error) {
	return s.repo.GetByID(id)
}

func (s *BrandService) GetAllBrands() ([]*Brand, error) {
	return s.repo.GetAll()
}

func (s *BrandService) CreateBrand(brand *Brand) (*Brand, error) {
	return s.repo.Create(brand)
}

func (s *BrandService) UpdateBrand(brand *Brand) (*Brand, error) {
	if brand.ID == 0 {
		return nil, errors.New("invalid brand ID")
	}
	newBrand, err := s.repo.Update(brand)
	if err != nil {
		return nil, err
	}

	return newBrand, nil
}

func (s *BrandService) DeleteBrand(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
