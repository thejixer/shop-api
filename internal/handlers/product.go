package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/shop-api/internal/data-process-layer"
	"github.com/thejixer/shop-api/internal/models"
)

func (h *HandlerService) CreateProduct(c echo.Context) error {
	body := models.CreateProductDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "provide sufficient data")
	}

	thisProduct, err := h.store.ProductRepo.Create(body)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToProductDto(thisProduct))

}

func (h *HandlerService) EditProduct(c echo.Context) error {

	id := c.Param("id")
	body := models.CreateProductDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "provide sufficient data")
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}

	var product *models.Product
	var er error

	product, er = h.store.ProductRepo.FindById(productId)
	if er != nil {
		return WriteReponse(c, http.StatusNotFound, er.Error())
	}

	product, er = h.store.ProductRepo.Edit(productId, body)
	if er != nil {
		return WriteReponse(c, http.StatusNotFound, er.Error())
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToProductDto(product))
}

func (h *HandlerService) GetProducts(c echo.Context) error {
	text := c.QueryParam("text")
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

	products, count, fErr := h.store.ProductRepo.Find(text, page, limit)
	if fErr != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}

	result := dataprocesslayer.ConvertToLLProductDto(products, count)

	return c.JSON(http.StatusOK, result)
}

func (h *HandlerService) GetProduct(c echo.Context) error {

	fmt.Println("im getting called")
	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	product, err := h.store.ProductRepo.FindById(productId)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToProductDto(product))

}
