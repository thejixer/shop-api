package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyRoutes(e *echo.Echo) {

	e.GET("/", s.handlerService.HandleHelloWorld)
	e.POST("/checkout", s.handlerService.CheckOut, s.handlerService.AuthGaurd)

	admin := e.Group("/admin")
	admin.POST("/me", s.handlerService.HandleAdminMe, s.handlerService.AdminGaurd)
	admin.POST("/create", s.handlerService.CreateAdmin, s.handlerService.AdminGaurd)
	admin.POST("/update-permissions", s.handlerService.UpdatePermissions, s.handlerService.AdminGaurd)

	user := e.Group("/user")
	user.GET("/", s.handlerService.GetUsers, s.handlerService.AdminGaurd)
	user.GET("/:id", s.handlerService.GetSingleUser, s.handlerService.AdminGaurd)
	user.POST("/charge-balance", s.handlerService.ChargeBalance, s.handlerService.AuthGaurd)

	cart := e.Group("/cart")
	cart.POST("/", s.handlerService.AddToCart, s.handlerService.AuthGaurd)
	cart.GET("/", s.handlerService.GetMyCart, s.handlerService.AuthGaurd)
	cart.DELETE("/:id", s.handlerService.RemoveFromCart, s.handlerService.AuthGaurd)

	address := e.Group("/address")
	address.POST("/", s.handlerService.CreateAddress, s.handlerService.AuthGaurd)
	address.GET("/", s.handlerService.GetMyAddresses, s.handlerService.AuthGaurd)
	address.DELETE("/:id", s.handlerService.DeleteAddress, s.handlerService.AuthGaurd)

	s.ApplyAuthRoutes(e)
	s.ApplyProductRoutes(e)
	s.ApplyOrderRoutes(e)

}
