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
}
