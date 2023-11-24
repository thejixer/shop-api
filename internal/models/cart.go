package models

import "github.com/shopspring/decimal"

type CartRepository interface {
	Add(userId, productId, quantity int) error
	FindUsersItems(userId int) ([]*CartItem, error)
	Remove(userId, productId int) error
}

type CartItem struct {
	Quantity           int             `json:"quantity"`
	ProductId          int             `json:"productId"`
	ProductTitle       string          `json:"productTitle"`
	ProductPrice       decimal.Decimal `json:"productPrice"`
	ProductQuantity    int             `json:"productQuantity"`
	ProductDescription string          `json:"productDescription"`
}

type CartItemDto struct {
	Product  ProductDto `json:"product"`
	Quantity int        `json:"quantity"`
}

type CartDto struct {
	User       UserDto       `json:"user"`
	Items      []CartItemDto `json:"items"`
	TotalPrice float64       `json:"totalPrice"`
}
