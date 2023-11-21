package handlers

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/thejixer/shop-api/internal/models"
)

func WriteReponse(c echo.Context, s int, m string) error {
	return c.JSON(s, models.ResponseDTO{Msg: m, StatusCode: s})
}

func CreateUUID() string {

	env := os.Getenv("ENVIROMENT")
	if env == "DEV" || env == "TEST" {
		return "1111"
	}

	return uuid.New().String()
}

func FindSingleUser(h *HandlerService, id int) (*models.User, error) {

	var thisUser *models.User
	var err error
	thisUser = h.redisStore.GetUser(id)

	if thisUser != nil {
		return thisUser, nil
	}

	thisUser, err = h.store.UserRepo.FindById(id)
	if err != nil {
		return nil, errors.New("not found")
	}

	go h.redisStore.CacheUser(thisUser)

	return thisUser, nil
}
