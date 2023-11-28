package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrdeRepository interface {
	Create(data Order, cartItems []*CartItem) error
	FindById(id int) (*Order, error)
	FindOrderItemsOfSingleOrder(orderId, userId int) ([]*OrderItem, error)
}

type Order struct {
	Id         int             `json:"id"`
	UserId     int             `json:"userId"`
	Status     string          `json:"status"`
	TotalPrice decimal.Decimal `json:"totalPrice"`
	CreatedAt  time.Time       `json:"createdAt"`
}

type OrderItem struct {
	OrderId            int             `json:"orderId"`
	Quantity           int             `json:"quantity"`
	ProductId          int             `json:"productId"`
	ProductTitle       string          `json:"productTitle"`
	ProductPrice       decimal.Decimal `json:"productPrice"`
	ProductQuantity    int             `json:"productQuantity"`
	ProductDescription string          `json:"productDescription"`
	PriceAtTheTime     decimal.Decimal `json:"priceAtTheTime"`
}

type OrderItemDto struct {
	Product        ProductDto `json:"product"`
	Quantity       int        `json:"quantity"`
	PriceAtTheTime float64    `json:"priceAtTheTime"`
}

type OrderDto struct {
	Id        int            `json:"id"`
	User      UserDto        `json:"user"`
	Status    string         `json:"status"`
	Items     []OrderItemDto `json:"items"`
	CreatedAt time.Time      `json:"createdAt"`
}
