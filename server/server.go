package server

import (
	"github.com/labstack/echo/v4"
	"github.com/thejixer/shop-api/handlers"
)

type APIServer struct {
	listenAddr     string
	handlerService *handlers.HandlerService
}

func NewAPIServer(listenAddr string, handlerService *handlers.HandlerService) *APIServer {

	return &APIServer{
		listenAddr:     listenAddr,
		handlerService: handlerService,
	}
}

func (s *APIServer) Run() {
	e := echo.New()

	s.ApplyMiddlewares(e)
	s.ApplyRoutes(e)

	e.Logger.Fatal(e.Start(s.listenAddr))
}
