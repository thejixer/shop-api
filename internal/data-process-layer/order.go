package dataprocesslayer

import (
	"github.com/shopspring/decimal"
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

func MakeOrder(o *models.Order, i []*models.OrderItem, u models.UserDto, a *models.AddressDto) *models.OrderDto {

	var items []models.OrderItemDto

	var temp decimal.Decimal

	for _, e := range i {
		items = append(items, *ConvertOrderItemtoDto(e))
		Quantity := decimal.NewFromInt(int64(e.Quantity))
		tp := Quantity.Mul(e.PriceAtTheTime)
		temp = temp.Add(tp)
	}
	totalPrice, _ := temp.Float64()

	return &models.OrderDto{
		Id:         o.Id,
		User:       u,
		Address:    *a,
		Status:     o.Status,
		Items:      items,
		CreatedAt:  o.CreatedAt,
		TotalPrice: totalPrice,
	}

}
