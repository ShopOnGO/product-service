package productVariant

import (
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

func HandleProductVariantEvent(msg []byte, key string, productVariantSvc *ProductVariantService) error {
	logger.Infof("Получено сообщение для варианта продукта: %s", string(msg))

	var base BaseProductVariantEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения варианта: %w", err)
	}

	eventHandlers := map[string]func([]byte, *ProductVariantService) error{
		"create": HandleCreateProductVariantEvent,
		// "update": HandleUpdateProductVariantEvent,
		// "delete": HandleDeleteProductVariantEvent,
	}

	handler, exists := eventHandlers[base.Action]
	if !exists {
		return fmt.Errorf("неизвестное действие для варианта продукта: %s", base.Action)
	}

	return handler(msg, productVariantSvc)
}

// HandleCreateProductVariantEvent обрабатывает создание варианта продукта
func HandleCreateProductVariantEvent(msg []byte, productVariantSvc *ProductVariantService) error {
	var base BaseProductVariantEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения варианта: %w", err)
	}

	event := base.ProductVariant
	logger.Infof(
        "Создание варианта (prod=%d) от seller %d: SKU=%q, Price=%s, Discount=%s, Stock=%d, IsActive=%t",
        base.ProductID, base.UserID,
        event.SKU, event.Price.String(), event.Discount.String(), event.Stock, event.IsActive,
    )

	// if err := productVariantSvc.CheckProductOwnership(base.ProductID, base.UserID); err != nil {
	// 	logger.Warnf("Проверка прав доступа отклонена: %v", err)
	// 	return err
	// }
	
	newProductVariant := &ProductVariant{
		ProductID:  base.ProductID,
		SKU:        event.SKU,
		Price:      event.Price,
		Discount:   event.Discount,
		Stock:      event.Stock,
		IsActive:   event.IsActive,
		Sizes:      event.Sizes,
		Colors:     event.Colors,
		Barcode:    event.Barcode,
		ImageURLs:  event.Images,
		MinOrder:   event.MinOrder,
		Dimensions: event.Dimensions,
	}

	created, err := productVariantSvc.CreateProductVariant(newProductVariant)
	if err != nil {
		logger.Errorf("Ошибка при создании варианта продукта: %v", err)
		return err
	}

	logger.Infof("Вариант продукта успешно создан: %+v", created)
	return nil
}