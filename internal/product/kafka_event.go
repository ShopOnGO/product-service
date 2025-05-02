package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

func HandleProductEvent(msg []byte, key string, productSvc *ProductService, kafkaProducer *kafkaService.KafkaService) error {
	logger.Infof("Получено сообщение: %s", string(msg))

	var base BaseProductEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	eventHandlers := map[string]func([]byte, *ProductService, *kafkaService.KafkaService) error{
		"create": HandleCreateProductEvent,
		"media-stored": HandleMediaEvent,
		// "update": HandleUpdateProductEvent,
		// "delete": HandleDeleteProductEvent,
	}

	handler, exists := eventHandlers[base.Action]
	if !exists {
		return fmt.Errorf("неизвестное действие для продукта: %s", base.Action)
	}

	return handler(msg, productSvc, kafkaProducer)
}

func HandleCreateProductEvent(msg []byte, productSvc *ProductService, kafkaProducer *kafkaService.KafkaService) error {
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
	}

	createdProduct, err := productSvc.CreateProduct(newProduct)
	if err != nil {
		logger.Errorf("Ошибка при создании отзыва: %v", err)
		return err
	}
	logger.Infof("Продукт успешно создан: %+v", createdProduct)

	eventForMedia := ProductCreatedEventForMedia{
		Action:    "create",
		ProductID: createdProduct.ID,
		ImageKeys: event.ImageKeys,
		VideoKeys: event.VideoKeys,
	}

	value, err := json.Marshal(eventForMedia)
	if err != nil {
		logger.Errorf("Ошибка сериализации события product-created: %v", err)
		return err
	}

	ctx := context.Background()
	if err := kafkaProducer.Produce(ctx, []byte("product-created"), value); err != nil {
		logger.Errorf("Ошибка отправки сообщения в Kafka: %v", err)
		return err
	}

	return nil
}


func HandleMediaEvent(msg []byte, productSvc *ProductService, kafkaProducer *kafkaService.KafkaService) error {
	var event MediaUpdateEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("ошибка десериализации события обновления медиа: %w", err)
	}

	logger.Infof("Обновление медиа для продукта ID %d: images=%v, video=%q",
		event.ProductID, event.ImageURLs, event.VideoURLs)

	if err := productSvc.UpdateProductMedia(event.ProductID, event.ImageURLs, event.VideoURLs); err != nil {
		logger.Errorf("Ошибка при обновлении медиа: %v", err)
		return err
	}

	logger.Infof("Медиа успешно обновлены для продукта ID %d", event.ProductID)
	return nil
}
