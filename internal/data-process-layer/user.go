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

func ConvertToLLUserDto(a []*models.User, count int) models.LL_UserDto {

	var users []models.UserDto

	for _, s := range a {
		users = append(users, ConvertToUserDto(s))
	}

	return models.LL_UserDto{
		Total:  count,
		Result: users,
	}

}

func ConvertToAdminDto(u *models.User) models.AdminDto {
	b, _ := u.Balance.Float64()
	return models.AdminDto{
		ID:          u.ID,
		Role:        u.Role,
		Name:        u.Name,
		Email:       u.Email,
		Balance:     b,
		Permissions: u.Permissions,
		CreatedAt:   u.CreatedAt,
	}
}

func ConvertToLLAdminDto(a []*models.User, count int) models.LL_AdminDto {

	var users []models.AdminDto

	for _, s := range a {
		users = append(users, ConvertToAdminDto(s))
	}

	return models.LL_AdminDto{
		Total:  count,
		Result: users,
	}

}
