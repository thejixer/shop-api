package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/shop-api/internal/data-process-layer"
	"github.com/thejixer/shop-api/internal/models"
)

func (h *HandlerService) CheckOut(c echo.Context) error {

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}
	cartItems, err := h.store.CartRepo.FindUsersItems(me.ID)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops. this one's on us")
	}

	if len(cartItems) == 0 {
		return WriteReponse(c, http.StatusBadRequest, "there is nothing in your cart")

	}

	user := dataprocesslayer.ConvertToUserDto(me)
	cart := dataprocesslayer.ConvertItemsToCart(user, cartItems)
	if user.Balance < cart.TotalPrice {
		return WriteReponse(c, http.StatusBadRequest, "you don't have enough credit, pleasse charge your account")
	}

	thisOrder := models.Order{
		UserId:    user.ID,
		Status:    "created",
		CreatedAt: time.Now().UTC(),
	}

	if err := h.store.OrderRepo.Create(thisOrder, cartItems); err != nil {
		if strings.Contains(err.Error(), "products_quantity_check") {
			return WriteReponse(c, http.StatusBadRequest, "our shop does not have that much of one of your requested products")
		}
		return WriteReponse(c, http.StatusInternalServerError, "oops. this one's on us")
	}

	return WriteReponse(c, http.StatusCreated, "successfully submited your order")
}

func (h *HandlerService) GetOrder(c echo.Context) error {
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	order, err := h.store.OrderRepo.FindById(orderId)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	if order.UserId != me.ID {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	orderItems, err := h.store.OrderRepo.FindOrderItemsOfSingleOrder(order.Id, me.ID)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops, this one's on us")
	}
	meDto := dataprocesslayer.ConvertToUserDto(me)

	thisOrder := dataprocesslayer.MakeOrder(order, orderItems, meDto)

	return c.JSON(http.StatusOK, thisOrder)
}
