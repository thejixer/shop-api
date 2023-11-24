package dataprocesslayer

import "github.com/thejixer/shop-api/internal/models"

func ConvertToProductDto(p *models.Product) models.ProductDto {
	b, _ := p.Price.Float64()
	return models.ProductDto{
		Id:          p.Id,
		Title:       p.Title,
		Price:       b,
		Quantity:    p.Quantity,
		Description: p.Description,
	}
}
