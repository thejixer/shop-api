package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyProductRoutes(e *echo.Echo) {

	product := e.Group("/product")
	product.POST("/", s.handlerService.CreateProduct, s.handlerService.AdminGaurd)
	product.POST("/:id", s.handlerService.EditProduct, s.handlerService.AdminGaurd)
	product.GET("/", s.handlerService.GetProducts)
	product.GET("/:id", s.handlerService.GetProduct)
}
