package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/shop-api/internal/data-process-layer"
	"github.com/thejixer/shop-api/internal/models"
	"github.com/thejixer/shop-api/internal/pdf"
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

	body := models.CheckOutDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "addressId is missing")
	}

	user := dataprocesslayer.ConvertToUserDto(me)
	cart := dataprocesslayer.ConvertItemsToCart(user, cartItems)

	if user.Balance < cart.TotalPrice {
		return WriteReponse(c, http.StatusBadRequest, "you don't have enough credit, pleasse charge your account")
	}

	thisAddress, err := h.store.AddressRepo.FindById(body.AddressId)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}

	if thisAddress.UserId != me.ID {
		return WriteReponse(c, http.StatusForbidden, "forbiden")
	}

	thisOrder := models.Order{
		UserId:    user.ID,
		Status:    "created",
		CreatedAt: time.Now().UTC(),
		AddressId: body.AddressId,
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

	if me.Role != "admin" && order.UserId != me.ID {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var makeOrderErr error
	thisOrder := new(models.OrderDto)
	go h.store.OrderRepo.MakeOrder(order, thisOrder, &makeOrderErr, &wg)
	wg.Wait()

	if makeOrderErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	return c.JSON(http.StatusOK, thisOrder)
}

func (h *HandlerService) GetMyOrders(c echo.Context) error {
	p := c.QueryParam("page")
	l := c.QueryParam("limit")

	var page int
	var limit int
	var err error

	page, err = strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err = strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	result, count, err := h.store.OrderRepo.GetOrdersByUserId(me.ID, page, limit)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	var orders = make([]models.OrderDto, len(result))
	var wg sync.WaitGroup

	for i, o := range result {
		wg.Add(1)
		var makeOrderErr error
		orders[i] = *new(models.OrderDto)
		go h.store.OrderRepo.MakeOrder(o, &orders[i], &makeOrderErr, &wg)
		if makeOrderErr != nil {
			return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
		}
	}
	wg.Wait()

	data := &models.LL_OrderDto{
		Total:  count,
		Result: orders,
	}

	return c.JSON(http.StatusOK, data)
}

func (h *HandlerService) AdminGetOrders(c echo.Context) error {

	p := c.QueryParam("page")
	l := c.QueryParam("limit")
	status := c.QueryParam("status")
	u := c.QueryParam("userId")

	var page int
	var limit int
	var err error

	page, err = strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err = strconv.Atoi(l)
	if err != nil {
		limit = 10
	}
	userId, err := strconv.Atoi(u)
	if err != nil {
		userId = 0
	}

	var meErr error
	_, meErr = GetMe(&c)
	if meErr != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	result, count, err := h.store.OrderRepo.QueryOrders(userId, status, page, limit)

	var orders = make([]models.OrderDto, len(result))
	var wg sync.WaitGroup

	for i, o := range result {
		wg.Add(1)
		var makeOrderErr error
		orders[i] = *new(models.OrderDto)
		go h.store.OrderRepo.MakeOrder(o, &orders[i], &makeOrderErr, &wg)
		if makeOrderErr != nil {
			return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
		}
	}
	wg.Wait()

	data := &models.LL_OrderDto{
		Total:  count,
		Result: orders,
	}

	return c.JSON(http.StatusOK, data)

}

func (h *HandlerService) VerifyOrder(c echo.Context) error {
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}
	if hasPermission := PermissionChecker(me, "backoffice"); !hasPermission {
		return WriteReponse(c, http.StatusForbidden, "forbidden resources")
	}
	order, err := h.store.OrderRepo.FindById(orderId)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	if order.Status != "created" {
		return WriteReponse(c, http.StatusBadRequest, "bad request")
	}
	updateErr := h.store.OrderRepo.SetStatus(orderId, "verified")
	if updateErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	return WriteReponse(c, http.StatusAccepted, "successfully updated the order")

}

func (h *HandlerService) PackageOrder(c echo.Context) error {
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}
	if hasPermission := PermissionChecker(me, "stock"); !hasPermission {
		return WriteReponse(c, http.StatusForbidden, "forbidden resources")
	}

	order, err := h.store.OrderRepo.FindById(orderId)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	if order.Status != "verified" {
		return WriteReponse(c, http.StatusBadRequest, "bad request")
	}

	updateErr := h.store.OrderRepo.SetStatus(orderId, "packaged")
	if updateErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	return WriteReponse(c, http.StatusAccepted, "successfully updated the order")

}

func (h *HandlerService) SendOrder(c echo.Context) error {
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}
	if hasPermission := PermissionChecker(me, "backoffice"); !hasPermission {
		return WriteReponse(c, http.StatusForbidden, "forbidden resources")
	}

	order, err := h.store.OrderRepo.FindById(orderId)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	if order.Status != "packaged" {
		return WriteReponse(c, http.StatusBadRequest, "bad request")
	}

	if err := h.redisStore.CreateShipment(order.Id, CustomUUID(6)); err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	updateErr := h.store.OrderRepo.SetStatus(orderId, "sent")
	if updateErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var makeOrderErr error
	thisOrder := new(models.OrderDto)
	go h.store.OrderRepo.MakeOrder(order, thisOrder, &makeOrderErr, &wg)
	wg.Wait()
	h.scheduelerService.SchedueleShipmentNotification(thisOrder)

	return WriteReponse(c, http.StatusAccepted, "successfully updated the order")

}

func (h *HandlerService) GetShipmentCode(c echo.Context) error {
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}
	order, err := h.store.OrderRepo.FindById(orderId)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	if order.UserId != me.ID {
		return WriteReponse(c, http.StatusForbidden, "forbidden")
	}

	code, err := h.redisStore.GetShipmentCode(order.Id)

	return c.JSON(http.StatusOK, code)
}

func (h *HandlerService) DeliverOrder(c echo.Context) error {

	body := models.DeliverOrderDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid data")
	}
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid data")
	}

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	if hasPermission := PermissionChecker(me, "shipper"); !hasPermission {
		return WriteReponse(c, http.StatusForbidden, "forbidden resources")
	}

	order, err := h.store.OrderRepo.FindById(orderId)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	if order.Status != "sent" {
		return WriteReponse(c, http.StatusBadRequest, "bad request")
	}

	code, err := h.redisStore.GetShipmentCode(order.Id)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops, this one's on us")
	}

	if code != body.Code {
		return WriteReponse(c, http.StatusForbidden, "forbidden resources")
	}

	updateErr := h.store.OrderRepo.SetStatus(orderId, "delivered")
	if updateErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	h.redisStore.DelShipment(order.Id)

	return WriteReponse(c, http.StatusAccepted, "successfully updated the order")

}

func (h *HandlerService) DownloadOrderPDF(c echo.Context) error {
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

	if me.Role != "admin" && order.UserId != me.ID {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var makeOrderErr error
	thisOrder := new(models.OrderDto)
	go h.store.OrderRepo.MakeOrder(order, thisOrder, &makeOrderErr, &wg)
	wg.Wait()

	if makeOrderErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	pdfAddress, err := pdf.GeneratePDF(thisOrder)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops this one's on us")
	}

	return c.File(pdfAddress)
}
