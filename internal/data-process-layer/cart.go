package dataprocesslayer

import (
	"github.com/shopspring/decimal"
	"github.com/thejixer/shop-api/internal/models"
)

func ConvertToCartItemDto(i *models.CartItem) models.CartItemDto {
	price, _ := i.ProductPrice.Float64()

	p := models.ProductDto{
		Id:          i.ProductId,
		Title:       i.ProductTitle,
		Price:       price,
		Quantity:    i.ProductQuantity,
		Description: i.ProductDescription,
	}

	return models.CartItemDto{
		Product:  p,
		Quantity: i.Quantity,
	}
}

func ConvertItemsToCart(user models.UserDto, items []*models.CartItem) models.CartDto {

	var arr []models.CartItemDto

	var temp decimal.Decimal

	for _, s := range items {
		arr = append(arr, ConvertToCartItemDto(s))

		Quantity := decimal.NewFromInt(int64(s.Quantity))
		tp := Quantity.Mul(s.ProductPrice)
		temp = temp.Add(tp)
	}

	TotalPrice, _ := temp.Float64()

	return models.CartDto{
		User:       user,
		Items:      arr,
		TotalPrice: TotalPrice,
	}
}
