package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/shop-api/internal/data-process-layer"
	"github.com/thejixer/shop-api/internal/models"
)

func (h *HandlerService) AddToCart(c echo.Context) error {

	body := models.AddtoCartDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid data")
	}

	thisProduct, err := h.store.ProductRepo.FindById(body.ProductId)
	if err != nil || thisProduct == nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid data")
	}

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	err2 := h.store.CartRepo.Add(me.ID, body.ProductId, body.Quantity)

	if err2 != nil {
		fmt.Println(err2)
		return WriteReponse(c, http.StatusInternalServerError, "oops, this one's on us")
	}

	return WriteReponse(c, http.StatusCreated, "successfully added to cart")
}

func (h *HandlerService) GetMyCart(c echo.Context) error {
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	items, err := h.store.CartRepo.FindUsersItems(me.ID)
	if err != nil {
		fmt.Println(err)
		return WriteReponse(c, http.StatusInternalServerError, "oops. this one's on us")
	}

	user := dataprocesslayer.ConvertToUserDto(me)
	cart := dataprocesslayer.ConvertItemsToCart(user, items)

	return c.JSON(http.StatusOK, cart)
}

func (h *HandlerService) RemoveFromCart(c echo.Context) error {
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid data")
	}
	err2 := h.store.CartRepo.Remove(me.ID, productId)
	if err2 != nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid data")
	}

	return WriteReponse(c, http.StatusAccepted, "removed successfully")
}
