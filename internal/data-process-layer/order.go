package dataprocesslayer

import "github.com/thejixer/shop-api/internal/models"

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

func MakeOrder(o *models.Order, i []*models.OrderItem, u models.UserDto) *models.OrderDto {

	var items []models.OrderItemDto

	for _, e := range i {
		items = append(items, *ConvertOrderItemtoDto(e))
	}

	return &models.OrderDto{
		Id:        o.Id,
		User:      u,
		Status:    o.Status,
		Items:     items,
		CreatedAt: o.CreatedAt,
	}

}
