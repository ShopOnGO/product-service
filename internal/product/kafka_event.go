package product

import (
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

func HandleProductEvent(msg []byte, key string, productSvc *ProductService) error {
	logger.Infof("Получено сообщение: %s", string(msg))

	var base BaseProductEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	eventHandlers := map[string]func([]byte, *ProductService) error{
		"create": HandleCreateProductEvent,
		// "update": HandleUpdateProductEvent,
		// "delete": HandleDeleteProductEvent,
	}

	handler, exists := eventHandlers[base.Action]
	if !exists {
		return fmt.Errorf("неизвестное действие для продукта: %s", base.Action)
	}

	return handler(msg, productSvc)
}

func HandleCreateProductEvent(msg []byte, productSvc *ProductService) error {
	var base BaseProductEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	event := base.Product
	logger.Infof("Получены данные для создания продукта: name=%q, category_id=%d, brand_id=%d, price=%d",
		event.Name, event.CategoryID, event.BrandID, event.Price)

	newProduct := &Product{
		Name:        event.Name,
		Description: event.Description,
		Price:       event.Price,
		Discount:    event.Discount,
		IsActive:    event.IsActive,
		CategoryID:  event.CategoryID,
		BrandID:     event.BrandID,
		Images:      event.Images,
		VideoURL:    event.VideoURL,
	}

	createdProduct, err := productSvc.CreateProduct(newProduct)
	if err != nil {
		logger.Errorf("Ошибка при создании отзыва: %v", err)
		return err
	}
	logger.Infof("Продукт успешно создан: %+v", createdProduct)
	return nil
}
