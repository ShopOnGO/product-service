package productVariant

import (
	"errors"

	"github.com/ShopOnGO/product-service/pkg/db"

	"gorm.io/gorm"
)

type ProductVariantRepository struct {
	Database *db.Db
}

func NewProductVariantRepository(database *db.Db) *ProductVariantRepository {
	return &ProductVariantRepository{
		Database: database,
	}
}

func (repo *ProductVariantRepository) Create(variant *ProductVariant) (*ProductVariant, error) {
	result := repo.Database.DB.Create(variant)
	if result.Error != nil {
		return nil, result.Error
	}
	return variant, nil
}

func (repo *ProductVariantRepository) GetBySKU(sku string) (*ProductVariant, error) {
	var variant ProductVariant
	result := repo.Database.DB.
		Where("sku = ?", sku).
		First(&variant)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &variant, result.Error
}

// GetByProductID возвращает все варианты для продукта
func (repo *ProductVariantRepository) GetByProductID(productID uint, includeInactive bool) ([]ProductVariant, error) {
	var variants []ProductVariant
	query := repo.Database.DB.Where("product_id = ?", productID)

	if !includeInactive {
		query = query.Where("is_active = true")
	}

	result := query.Find(&variants)
	return variants, result.Error
}

// GetActive возвращает только активные варианты
func (repo *ProductVariantRepository) GetActive() ([]ProductVariant, error) {
	var variants []ProductVariant
	result := repo.Database.DB.
		Where("is_active = true").
		Find(&variants)
	return variants, result.Error
}

// GetByID возвращает вариант по его ID
func (repo *ProductVariantRepository) GetByID(id uint) (*ProductVariant, error) {
	var variant ProductVariant
	result := repo.Database.DB.
		Where("id = ?", id).
		First(&variant)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &variant, result.Error
}

func (repo *ProductVariantRepository) GetVariantsByIDs(ids []uint) ([]ProductVariant, error) {
    var variants []ProductVariant
    if len(ids) == 0 {
        return variants, nil
    }
    if err := repo.Database.DB.
        Where("id IN ?", ids).
        Find(&variants).
        Error; err != nil {
        return nil, err
    }
    return variants, nil
}

// GetByBarcode поиск по штрихкоду
func (repo *ProductVariantRepository) GetByBarcode(barcode string) (*ProductVariant, error) {
	var variant ProductVariant
	result := repo.Database.DB.
		Where("barcode = ?", barcode).
		First(&variant)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &variant, result.Error
}

// Update обновляет вариант продукта
func (repo *ProductVariantRepository) Update(variant *ProductVariant) (*ProductVariant, error) {
	result := repo.Database.DB.Model(&ProductVariant{}).
	Where("id = ?", variant.ID).
	Updates(variant)

	if result.Error != nil {
		return nil, result.Error
	}
	return variant, nil
}

// SoftDelete мягкое удаление
func (repo *ProductVariantRepository) SoftDelete(id uint) error {
	return repo.Database.DB.Delete(&ProductVariant{}, id).Error
}

// GetByFilters поиск с фильтрами
func (repo *ProductVariantRepository) GetByFilters(filters map[string]interface{}, limit, offset int) ([]ProductVariant, error) {
	var variants []ProductVariant
	query := repo.Database.DB.Model(&ProductVariant{})

	for key, value := range filters {
		switch key {
		case "min_price":
			query = query.Where("price >= ?", value)
		case "max_price":
			query = query.Where("price <= ?", value)
		case "sizes":
			query = query.Where("JSON_CONTAINS(sizes, ?)", value)
		case "colors":
			query = query.Where("JSON_CONTAINS(colors, ?)", value)
		case "material":
			query = query.Where("material = ?", value)
		}
	}

	result := query.Limit(limit).Offset(offset).Find(&variants)
	return variants, result.Error
}



// GetAvailableStock возвращает доступное количество (Stock - ReservedStock)
func (repo *ProductVariantRepository) GetAvailableStock(variantID uint) (uint32, error) {
	var available struct {
		Available uint32
	}

	result := repo.Database.DB.Model(&ProductVariant{}).
		Select("(stock - reserved_stock) as available").
		Where("id = ?", variantID).
		Scan(&available)

	return available.Available, result.Error
}

// UpdateStock обновляет общий остаток на складе
func (repo *ProductVariantRepository) UpdateStock(variantID uint, newStock uint32) error {
	return repo.Database.DB.Model(&ProductVariant{}).
		Where("id = ?", variantID).
		Update("stock", newStock).Error
}

// ReserveStock резервирует указанное количество товара
func (repo *ProductVariantRepository) ReserveStock(variantID uint, quantity uint32) error {
	return repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		var variant ProductVariant
		if err := tx.First(&variant, variantID).Error; err != nil {
			return err
		}

		if variant.Stock < variant.ReservedStock+quantity {
			return errors.New("not enough stock")
		}

		return tx.Model(&variant).
			Update("reserved_stock", variant.ReservedStock+quantity).Error
	})
}

// ReleaseStock освобождает зарезервированный товар
func (repo *ProductVariantRepository) ReleaseStock(variantID uint, quantity uint32) error {
	return repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		var variant ProductVariant
		if err := tx.First(&variant, variantID).Error; err != nil {
			return err
		}
		if quantity > variant.ReservedStock {
			return errors.New("release quantity exceeds reserved stock")
		}
		newReserved := variant.ReservedStock - quantity
		
		return tx.Model(&variant).
			Update("reserved_stock", newReserved).Error
	})
}

// BulkUpdateStock массовое обновление стока
func (repo *ProductVariantRepository) BulkUpdateStock(variantStocks map[uint]uint32) error {
	return repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		for variantID, stock := range variantStocks {
			if err := tx.Model(&ProductVariant{}).
				Where("id = ?", variantID).
				Update("stock", stock).Error; err != nil {
				return err
			}
		}
		return nil
	})
}


