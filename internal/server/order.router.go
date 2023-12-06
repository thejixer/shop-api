package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyOrderRoutes(e *echo.Echo) {
	order := e.Group("/order")
	order.GET("/single/:id", s.handlerService.GetOrder, s.handlerService.Gaurd)
	order.GET("/my-orders", s.handlerService.GetMyOrders, s.handlerService.AuthGaurd)
	order.GET("/", s.handlerService.AdminGetOrders, s.handlerService.AdminGaurd)
	order.POST("/verify/:id", s.handlerService.VerifyOrder, s.handlerService.AdminGaurd)
	order.POST("/package/:id", s.handlerService.PackageOrder, s.handlerService.AdminGaurd)
	order.POST("/send/:id", s.handlerService.SendOrder, s.handlerService.AdminGaurd)
	order.GET("/shipment-code/:id", s.handlerService.GetShipmentCode, s.handlerService.AuthGaurd)
	order.POST("/deliver/:id", s.handlerService.DeliverOrder, s.handlerService.AdminGaurd)
	order.GET("/pdf/:id", s.handlerService.DownloadOrderPDF, s.handlerService.Gaurd)
}
