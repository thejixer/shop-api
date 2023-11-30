package dataprocesslayer

import "github.com/thejixer/shop-api/internal/models"

func ConvertToAddressDto(a *models.Address, u models.UserDto) *models.AddressDto {

	coordinates := models.PointDto{
		Lat: a.Lat,
		Lon: a.Lon,
	}

	return &models.AddressDto{
		Id:            a.Id,
		User:          u,
		Title:         a.Title,
		Coordinates:   coordinates,
		Address:       a.Address,
		RecieverName:  a.RecieverName,
		RecieverPhone: a.RecieverPhone,
	}
}

func ConvertToAddressDtos(a []*models.Address, u models.UserDto) []*models.AddressDto {

	arr := make([]*models.AddressDto, 0)

	for _, e := range a {
		item := ConvertToAddressDto(e, u)
		arr = append(arr, item)
	}

	return arr

}
