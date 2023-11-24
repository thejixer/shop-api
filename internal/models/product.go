package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProductRepository interface {
	Create(data CreateProductDto) (*Product, error)
	Edit(id int, data CreateProductDto) (*Product, error)
	FindById(id int) (*Product, error)
	Find(text string, page, limit int) ([]*Product, error)
}

type Product struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"createdAt"`
}

type ProductDto struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
}
