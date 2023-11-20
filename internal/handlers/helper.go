package handlers

import (
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/thejixer/shop-api/internal/models"
)

func WriteReponse(c *echo.Context, s int, m string) error {
	return (*c).JSON(s, models.ResponseDTO{Msg: m, StatusCode: s})
}

func CreateUUID() string {

	env := os.Getenv("ENVIROMENT")
	if env == "DEV" || env == "TEST" {
		return "1111"
	}

	return uuid.New().String()
}
