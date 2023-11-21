package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyAuthRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	auth.POST("/signup", s.handlerService.HangeSingup)
	auth.POST("/request-verificationCode", s.handlerService.HandleRequestVerificationEmail)
	auth.GET("/verify-email", s.handlerService.HandleEmailVerification)
	auth.POST("/login", s.handlerService.HandleLogin)
	auth.POST("/me", s.handlerService.HandleMe, s.handlerService.AuthGaurd)
	auth.POST("/request-change-password", s.handlerService.HandleRequestChangePassword)
	auth.GET("/verify-changepassword-request", s.handlerService.HandleVerifyChangePasswordRequest)
	auth.POST("/change-password", s.handlerService.HandleChangePassword)
}
