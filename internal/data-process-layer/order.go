package dataprocesslayer

import (
	"github.com/thejixer/shop-api/internal/models"
)

func ConvertOrderItemtoDto(i *models.OrderItem) *models.OrderItemDto {
	currentPrice, _ := i.ProductPrice.Float64()
	price, _ := i.PriceAtTheTime.Float64()

	p := models.ProductDto{
		Id:          i.ProductId,
		Title:       i.ProductTitle,
		Price:       currentPrice,
		Quantity:    i.ProductQuantity,
		Description: i.ProductDescription,
	}

	return &models.OrderItemDto{
		Product:        p,
		Quantity:       i.Quantity,
		PriceAtTheTime: price,
	}
}
