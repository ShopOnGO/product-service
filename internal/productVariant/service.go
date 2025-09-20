package productVariant

import (
	"errors"
	"fmt"

	// "github.com/ShopOnGO/product-service/pkg/interfaces"
)

type ProductVariantService struct {
	repo *ProductVariantRepository
	// productRepo *interfaces.ProductChecker
}

func NewProductVariantService(repo *ProductVariantRepository) *ProductVariantService {
	return &ProductVariantService{
		repo: repo,
		// productRepo: productRepo,
	}
}

func (s *ProductVariantService) CreateProductVariant(variant *ProductVariant) (*ProductVariant, error) {
	if variant.SKU == "" {
		return nil, errors.New("SKU is required")
	}
	// проверка что такой ID продукта есть
	// exists, err := s.productRepo.ExistsByID(variant.ProductID)
	// if !exists || err != nil {
	// 	return nil, err
	// }
	
	// Проверка уникальности SKU
	existing, err := s.repo.GetBySKU(variant.SKU)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("product variant with SKU %s already exists", variant.SKU)
	}
	// Дополнительные проверки могут быть добавлены здесь (например, валидация размеров, цветов и пр.)
	return s.repo.Create(variant)
}

func (s *ProductVariantService) GetProductVariantByID(id uint) (*ProductVariant, error) {
	if id == 0 {
		return nil, errors.New("invalid product variant ID")
	}
	return s.repo.GetVariantByID(id)
}

func (s *ProductVariantService) GetVariantsByIDs(ids []uint) ([]ProductVariant, error) {
    return s.repo.GetVariantsByIDs(ids)
}

func (s *ProductVariantService) UpdateProductVariantByInput(variantID uint, input UpdateProductVariantPayload) (*ProductVariant, error) {
	if variantID == 0 {
		return nil, errors.New("variant ID is required for update")
	}
	// Получаем существующий вариант
	existing, err := s.repo.GetVariantByID(variantID)
	if err != nil {
		return nil, fmt.Errorf("product variant with ID %d not found: %w", variantID, err)
	}
	if existing == nil {
		return nil, fmt.Errorf("product variant with ID %d not found", variantID)
	}

	// Обновляем поля, если входные данные заданы
	if input.Price != nil {
		existing.Price = *input.Price
	}
	if input.Discount != nil {
		existing.Discount = *input.Discount
	}
	if input.ReservedStock != nil {
		existing.ReservedStock = *input.ReservedStock
	}
	if input.Sizes != nil {
		existing.Sizes = *input.Sizes
	}
	if input.Colors != nil {
		existing.Colors = *input.Colors
	}
	if input.Stock != nil {
		existing.Stock = *input.Stock
	}
	if input.Barcode != nil {
		existing.Barcode = *input.Barcode
	}
	if input.IsActive != nil {
		existing.IsActive = *input.IsActive
	}
	if input.ImageURLs != nil {
		existing.ImageURLs = *input.ImageURLs
	}
	if input.MinOrder != nil {
		existing.MinOrder = *input.MinOrder
	}
	if input.Dimensions != nil {
		existing.Dimensions = *input.Dimensions
	}

	// Вызываем репозиторий для обновления
	return s.repo.Update(existing)
}

// DeleteProductVariant выполняет мягкое удаление варианта продукта.
func (s *ProductVariantService) DeleteProductVariant(id uint) error {
	if id == 0 {
		return errors.New("invalid product variant ID")
	}
	return s.repo.SoftDelete(id)
}

// ReserveStock резервирует указанное количество товара, если доступно.
func (s *ProductVariantService) ReserveStock(variantID uint, quantity uint32) error {
	if quantity == 0 {
		return errors.New("reserve quantity must be greater than zero")
	}
	return s.repo.ReserveStock(variantID, quantity)
}

// ReleaseStock освобождает указанное количество зарезервированного товара.
// Добавлена базовая проверка, чтобы не освободить больше, чем зарезервировано.
func (s *ProductVariantService) ReleaseStock(variantID uint, quantity uint32) error {
	if quantity == 0 {
		return errors.New("release quantity must be greater than zero")
	}
	// Дополнительная логика: проверка, чтобы не произошло переполнение (underflow)
	variant, err := s.repo.GetVariantByID(variantID)
	if err != nil {
		return err
	}
	if quantity > variant.ReservedStock {
		return errors.New("release quantity exceeds reserved stock")
	}
	return s.repo.ReleaseStock(variantID, quantity)
}

// UpdateStock обновляет общее количество товара для варианта.
func (s *ProductVariantService) UpdateStock(variantID uint, newStock uint32) error {
	return s.repo.UpdateStock(variantID, newStock)
}

// GetAvailableStock возвращает доступное количество товара (stock - reserved_stock).
func (s *ProductVariantService) GetAvailableStock(variantID uint) (uint32, error) {
	return s.repo.GetAvailableStock(variantID)
}

// GetBySKU возвращает вариант продукта по артикулу.
func (s *ProductVariantService) GetBySKU(sku string) (*ProductVariant, error) {
	if sku == "" {
		return nil, errors.New("SKU is required")
	}
	return s.repo.GetBySKU(sku)
}


// func (s *ProductService) CheckProductOwnership(productID, userID uint) error {
//     owned, err := s.productRepo.IsProductOwnedByUser(productID, userID)
//     if err != nil {
//         return fmt.Errorf("ошибка при проверке прав владения продуктом: %w", err)
//     }

//     if !owned {
//         return fmt.Errorf("forbidden: product %d does not belong to user %d", productID, userID)
//     }

//     return nil
// }
