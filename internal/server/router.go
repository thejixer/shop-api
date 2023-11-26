package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyRoutes(e *echo.Echo) {

	e.GET("/", s.handlerService.HandleHelloWorld)

	s.ApplyAuthRoutes(e)

	admin := e.Group("/admin")
	admin.POST("/me", s.handlerService.HandleMe, s.handlerService.AdminGaurd)
	admin.POST("/create", s.handlerService.CreateAdmin, s.handlerService.AdminGaurd)

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
}
