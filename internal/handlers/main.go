package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thejixer/shop-api/internal/database"
	"github.com/thejixer/shop-api/internal/mailer"
	"github.com/thejixer/shop-api/internal/redis"
)

type HandlerService struct {
	store         *database.PostgresStore
	redisStore    *redis.RedisStore
	mailerService *mailer.MailerService
}

func NewHandlerService(store *database.PostgresStore, redisStore *redis.RedisStore, mailerService *mailer.MailerService) *HandlerService {
	return &HandlerService{
		store:         store,
		redisStore:    redisStore,
		mailerService: mailerService,
	}
}

func (h *HandlerService) HandleHelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World from shop-api!")
}
