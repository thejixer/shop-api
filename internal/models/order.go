package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrdeRepository interface {
	Create(data Order) (*Order, error)
}

type Order struct {
	Id        int       `json:"id"`
	UserId    int       `json:"userId"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	Quantity           int             `json:"quantity"`
	ProductId          int             `json:"productId"`
	ProductTitle       string          `json:"productTitle"`
	ProductPrice       decimal.Decimal `json:"productPrice"`
	ProductQuantity    int             `json:"productQuantity"`
	ProductDescription string          `json:"productDescription"`
}

type OrderItemDto struct {
	Product  ProductDto `json:"product"`
	Quantity int        `json:"quantity"`
}

type OrderDto struct {
	Id        int            `json:"id"`
	User      UserDto        `json:"user"`
	Status    string         `json:"status"`
	Items     []OrderItemDto `json:"items"`
	CreatedAt time.Time      `json:"createdAt"`
}
