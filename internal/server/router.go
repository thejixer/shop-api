package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyRoutes(e *echo.Echo) {

	e.GET("/", s.handlerService.HandleHelloWorld)

	s.ApplyAuthRoutes(e)

	admin := e.Group("/admin")
	admin.POST("/me", s.handlerService.HandleAdminMe, s.handlerService.AdminGaurd)
	admin.POST("/create", s.handlerService.CreateAdmin, s.handlerService.AdminGaurd)
	admin.POST("/update-permissions", s.handlerService.UpdatePermissions, s.handlerService.AdminGaurd)

	user := e.Group("/user")
	user.GET("/", s.handlerService.GetUsers, s.handlerService.AdminGaurd)
	user.GET("/:id", s.handlerService.GetSingleUser, s.handlerService.AdminGaurd)
	user.POST("/charge-balance", s.handlerService.ChargeBalance, s.handlerService.AuthGaurd)

	product := e.Group("/product")
	product.POST("/", s.handlerService.CreateProduct, s.handlerService.AdminGaurd)
	product.POST("/:id", s.handlerService.EditProduct, s.handlerService.AdminGaurd)
	product.GET("/", s.handlerService.GetProducts)
	product.GET("/:id", s.handlerService.GetProduct)

	cart := e.Group("/cart")
	cart.POST("/", s.handlerService.AddToCart, s.handlerService.AuthGaurd)
	cart.GET("/", s.handlerService.GetMyCart, s.handlerService.AuthGaurd)
	cart.DELETE("/:id", s.handlerService.RemoveFromCart, s.handlerService.AuthGaurd)

	e.POST("/checkout", s.handlerService.CheckOut, s.handlerService.AuthGaurd)

	address := e.Group("/address")
	address.POST("/", s.handlerService.CreateAddress, s.handlerService.AuthGaurd)
	address.GET("/", s.handlerService.GetMyAddresses, s.handlerService.AuthGaurd)
	address.DELETE("/:id", s.handlerService.DeleteAddress, s.handlerService.AuthGaurd)

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
