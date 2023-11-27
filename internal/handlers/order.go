package handlers

import (
	"net/http"
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

	if _, err := h.store.OrderRepo.Create(thisOrder); err != nil {
		if strings.Contains(err.Error(), "your cart") {
			return WriteReponse(c, http.StatusBadRequest, err.Error())
		}
		if strings.Contains(err.Error(), "products_quantity_check") {
			return WriteReponse(c, http.StatusBadRequest, "our shop does not have that much of one of your requested products")
		}
		return WriteReponse(c, http.StatusInternalServerError, "oops. this one's on us")
	}

	return WriteReponse(c, http.StatusCreated, "successfully submited your order")
}
