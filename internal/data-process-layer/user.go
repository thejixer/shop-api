package dataprocesslayer

import (
	"github.com/thejixer/shop-api/internal/models"
)

func ConvertToUserDto(u *models.User) models.UserDto {
	b, _ := u.Balance.Float64()
	return models.UserDto{
		ID:        u.ID,
		Role:      u.Role,
		Name:      u.Name,
		Email:     u.Email,
		Balance:   b,
		CreatedAt: u.CreatedAt,
	}
}
