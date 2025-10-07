package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/ShopOnGO/product-service/internal/productVariant"
)

func HandleProductEvent(msg []byte, key string, productSvc *ProductService, productVariantSvc *productVariant.ProductVariantService, kafkaProducer *kafkaService.KafkaService) error {
	logger.Infof("Получено сообщение: %s", string(msg))

	var base BaseProductEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	eventHandlers := map[string]func([]byte, *ProductService, *productVariant.ProductVariantService, *kafkaService.KafkaService) error{
		"create": HandleCreateProductEvent,
		"media-stored": HandleMediaEvent,
		// "update": HandleUpdateProductEvent,
		// "delete": HandleDeleteProductEvent,
	}

	handler, exists := eventHandlers[base.Action]
	if !exists {
		return fmt.Errorf("неизвестное действие для продукта: %s", base.Action)
	}

	return handler(msg, productSvc, productVariantSvc, kafkaProducer)
}

func HandleCreateProductEvent(msg []byte, productSvc *ProductService, productVariantSvc *productVariant.ProductVariantService, kafkaProducer *kafkaService.KafkaService) error {
	var base BaseProductEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	event := base.Product
	logger.Infof("Получены данные для создания продукта: name=%q, category_id=%d, brand_id=%d",
		event.Name, event.CategoryID, event.BrandID)

	newProduct := &Product{
		Name:        event.Name,
		Description: event.Description,
		Material:    event.Material,
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

	var createdVariants []productVariant.ProductVariant

	for _, variantReq := range event.Variants {
		variant := &productVariant.ProductVariant{
			ProductID: createdProduct.ID,
			SKU:       variantReq.SKU,
			Price:     variantReq.Price,
			Discount:  variantReq.Discount,
			Sizes:     variantReq.Sizes,
			Colors:    variantReq.Colors,
			Stock:     variantReq.Stock,
			IsActive:  true,
    	}

		createdVariant, err := productVariantSvc.CreateProductVariant(variant)
		if err != nil {
			logger.Errorf("Ошибка при создании варианта: %v", err)
			return err
		}
		createdVariants = append(createdVariants, *createdVariant)
	}
	var variantsForEvent []*ProductVariantForEvent
	for _, v := range createdVariants {
		variantsForEvent = append(variantsForEvent, ConvertVariantToEvent(&v))
	}

	eventForMediaAndSearch := ProductCreatedEventForMediaAndSearch{
		Action:    		"create",
		ProductID: 		createdProduct.ID,
		Name:        	event.Name,
		Description: 	event.Description,
		Material:    	event.Material,
		Rating:      	event.Rating.InexactFloat64(),
		ReviewCount: 	event.ReviewCount,
		IsActive:    	event.IsActive,
		CategoryID:  	event.CategoryID,
		BrandID:     	event.BrandID,
		ImageKeys: 		event.ImageKeys,
		VideoKeys: 		event.VideoKeys,
		Variants:   	variantsForEvent,
	}

	value, err := json.Marshal(eventForMediaAndSearch)
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


func HandleMediaEvent(msg []byte, productSvc *ProductService, productVariantSvc *productVariant.ProductVariantService, kafkaProducer *kafkaService.KafkaService) error {
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


func ConvertVariantToEvent(v *productVariant.ProductVariant) *ProductVariantForEvent {
    return &ProductVariantForEvent{
        VariantID:     v.ID,
        SKU:           v.SKU,
        Price:         v.Price.InexactFloat64(),
        Discount:      v.Discount.InexactFloat64(),
        Sizes:         v.Sizes,
        Colors:        v.Colors,
        Stock:         v.Stock,
        Barcode:       v.Barcode,
        Dimensions:    v.Dimensions,
        ImageURLs:     v.ImageURLs, // pq.StringArray → []string
        MinOrder:      v.MinOrder,
        IsActive:      v.IsActive,
        ReservedStock: v.ReservedStock,
    }
}