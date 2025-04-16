package category

import (
	"errors"
	"fmt"
)

type CategoryService struct {
	repo *CategoryRepository
}

func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) CreateCategory(category *Category) (*Category, error) {
	existing, _ := s.repo.GetByName(category.Name)
	if existing != nil {
		return nil, errors.New("категория с таким именем уже существует")
	}

	if category.ParentCategoryID != nil {
		parent, err := s.repo.GetByID(*category.ParentCategoryID)
		if err != nil || parent == nil {
			return nil, errors.New("указана несуществующая родительская категория")
		}
	}
	return s.repo.Create(category)
}

func (s *CategoryService) GetFeaturedCategories(amount int) ([]Category, error) {
	return s.repo.GetFeaturedCategories(amount)
}

func (s *CategoryService) GetCategoryByName(name string) (*Category, error) {
	return s.repo.GetByName(name)
}

func (s *CategoryService) GetCategoryByID(id uint) (*Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) UpdateCategory(category *Category) (*Category, error) {
	existing, err := s.repo.GetByID(category.ID)
	if err != nil {
		return nil, fmt.Errorf("категория с ID %d не найдена: %w", category.ID, err)
	}

	if category.ParentCategoryID != nil && *category.ParentCategoryID == category.ID {
		return nil, errors.New("нельзя установить родительскую категорию саму на себя")
	}

	if category.Name != "" && category.Name != existing.Name {
		catWithSameName, err := s.repo.GetByName(category.Name)
		if err == nil && catWithSameName != nil && catWithSameName.ID != category.ID {
			return nil, fmt.Errorf("категория с именем %q уже существует", category.Name)
		}
	}

	if category.ParentCategoryID != nil {
		parent, err := s.repo.GetByID(*category.ParentCategoryID)
		if err != nil || parent == nil {
			return nil, errors.New("указанная родительская категория не существует")
		}
	}

	return s.repo.Update(category)
}

func (s *CategoryService) DeleteCategory(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("категория с ID %d не найдена: %w", id, err)
	}

	var subCategories []Category
	err = s.repo.Db.DB.Where("parent_category_id = ?", id).Find(&subCategories).Error
	if err != nil {
		return fmt.Errorf("ошибка при проверке подкатегорий: %w", err)
	}

	if len(subCategories) > 0 {
		return errors.New("нельзя удалить категорию, у которой есть подкатегории")
	}

	return s.repo.Delete(id)
}

