package models

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type OrdeRepository interface {
	Create(data Order, cartItems []*CartItem) error
	FindById(id int) (*Order, error)
	MakeOrder(order *Order, t *OrderDto, e *error, wg *sync.WaitGroup)
	GetOrdersByUserId(userId, page, limit int) ([]*Order, int, error)
	QueryOrders(userId int, status string, page, limit int) ([]*Order, int, error)
}

type Order struct {
	Id         int             `json:"id"`
	UserId     int             `json:"userId"`
	AddressId  int             `json:"addressId"`
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
	Id         int            `json:"id"`
	User       UserDto        `json:"user"`
	Address    AddressDto     `json:"address"`
	Status     string         `json:"status"`
	Items      []OrderItemDto `json:"items"`
	CreatedAt  time.Time      `json:"createdAt"`
	TotalPrice float64        `json:"totalPrice"`
}

type LL_OrderDto struct {
	Total  int        `json:"Total"`
	Result []OrderDto `json:"result"`
}
