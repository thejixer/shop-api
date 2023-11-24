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

func ConvertToLLProductDto(a []*models.Product, count int) models.LL_ProductDto {

	var users []models.ProductDto

	for _, s := range a {
		users = append(users, ConvertToProductDto(s))
	}

	return models.LL_ProductDto{
		Total:  count,
		Result: users,
	}

}
