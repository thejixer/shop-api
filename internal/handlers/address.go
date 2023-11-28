package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/shop-api/internal/data-process-layer"
	"github.com/thejixer/shop-api/internal/models"
)

func (h *HandlerService) CreateAddress(c echo.Context) error {
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	body := models.CreateAddressDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid data")
	}

	thisAddress := &models.Address{
		UserId:        me.ID,
		Title:         body.Title,
		Lon:           body.Lon,
		Lat:           body.Lat,
		Address:       body.Address,
		RecieverName:  body.RecieverName,
		RecieverPhone: body.RecieverPhone,
	}

	var insertErr error
	thisAddress, insertErr = h.store.AddressRepo.Create(thisAddress)

	if insertErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops, this one's on us")
	}

	user := dataprocesslayer.ConvertToUserDto(me)
	address := dataprocesslayer.ConvertToAddressDto(thisAddress, user)

	return c.JSON(http.StatusCreated, address)

}

func (h *HandlerService) GetMyAddresses(c echo.Context) error {
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	res, err := h.store.AddressRepo.FindByUserId(me.ID)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	var addresses []*models.Address
	for _, e := range res {
		fmt.Println(e.Deleted)
		if !e.Deleted {
			addresses = append(addresses, e)
		}
	}

	user := dataprocesslayer.ConvertToUserDto(me)
	myaddresses := dataprocesslayer.ConvertToAddressDtos(addresses, user)

	return c.JSON(http.StatusOK, myaddresses)

}

func (h *HandlerService) DeleteAddress(c echo.Context) error {
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}
	id := c.Param("id")
	addressId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	thisAddress, err := h.store.AddressRepo.FindById(addressId)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	if thisAddress.UserId != me.ID {
		return WriteReponse(c, http.StatusForbidden, "forbidden")
	}

	deletionErr := h.store.AddressRepo.Delete(thisAddress.Id)
	if deletionErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops, this one's on us")
	}

	return WriteReponse(c, http.StatusAccepted, "successfully deleted the address")

}
