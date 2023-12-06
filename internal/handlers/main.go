package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thejixer/shop-api/internal/database"
	"github.com/thejixer/shop-api/internal/mailer"
	"github.com/thejixer/shop-api/internal/redis"
	"github.com/thejixer/shop-api/internal/schedueler"
)

type HandlerService struct {
	store             *database.PostgresStore
	redisStore        *redis.RedisStore
	mailerService     *mailer.MailerService
	scheduelerService *schedueler.ScheduelerService
}

func NewHandlerService(
	store *database.PostgresStore,
	redisStore *redis.RedisStore,
	mailerService *mailer.MailerService,
	scheduelerService *schedueler.ScheduelerService,
) *HandlerService {
	return &HandlerService{
		store:             store,
		redisStore:        redisStore,
		mailerService:     mailerService,
		scheduelerService: scheduelerService,
	}
}

func (h *HandlerService) HandleHelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World from shop-api!")
}
