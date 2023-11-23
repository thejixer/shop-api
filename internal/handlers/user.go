package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/shop-api/internal/data-process-layer"
	"github.com/thejixer/shop-api/internal/models"
	"github.com/thejixer/shop-api/internal/utils"
	"github.com/thejixer/shop-api/pkg/encryption"
)

type CustomContext struct {
	echo.Context
	User *models.User
}

func (h *HandlerService) HangeSingup(c echo.Context) error {

	body := models.SignUpDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "lack of data")
	}

	thisUser, _ := h.store.UserRepo.FindByEmail(body.Email)

	if thisUser != nil {
		return WriteReponse(c, http.StatusBadRequest, "this email already exists in the database")
	}

	var err error
	thisUser, err = h.store.UserRepo.Create(body.Name, body.Email, body.Password, "user", false)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, err.Error())
	}

	verificationCode := CreateUUID()

	redisErr := h.redisStore.SetEmailVerificationCode(thisUser.Email, verificationCode)
	if redisErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	go h.mailerService.SendVerificationEmail(thisUser, verificationCode)

	return WriteReponse(c, http.StatusOK, "please check your email to verify your email")
}

func (h *HandlerService) HandleRequestVerificationEmail(c echo.Context) error {
	body := models.RequestVerificationEmailDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "lack of data")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "no user found")
	}

	if thisUser.IsEmailVerified {
		return WriteReponse(c, http.StatusBadRequest, "your email has already been verified")
	}

	verificationCode := CreateUUID()
	redisErr := h.redisStore.SetEmailVerificationCode(thisUser.Email, verificationCode)
	if redisErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	go h.mailerService.SendVerificationEmail(thisUser, verificationCode)

	return WriteReponse(c, http.StatusOK, "please check your email to verify your email")
}

func (h *HandlerService) HandleEmailVerification(c echo.Context) error {
	email := c.QueryParam("email")
	verificationCode := c.QueryParam("code")

	if email == "" || verificationCode == "" {
		return WriteReponse(c, http.StatusBadRequest, "insufficient data")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "user not found")
	}

	if val, err := h.redisStore.GetEmailVerificationCode(thisUser.Email); err != nil || val != verificationCode {
		return WriteReponse(c, http.StatusBadRequest, "code doesnt match")
	}

	updateErr := h.store.UserRepo.VerifyEmail(email)
	go h.redisStore.DeleteEmailVerificationCode(email)
	if updateErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	tokenString, err := utils.SignToken(thisUser.ID)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, models.TokenDTO{Token: tokenString})
}

func (h *HandlerService) HandleLogin(c echo.Context) error {
	body := models.LoginDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "lack of data")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(body.Email)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "bad requesst")
	}

	if !thisUser.IsEmailVerified {
		return WriteReponse(c, http.StatusUnauthorized, "your email is not verified")
	}

	if match := encryption.CheckPasswordHash(body.Password, thisUser.Password); !match {
		return WriteReponse(c, http.StatusBadRequest, "password doesnt match")
	}

	tokenString, err := utils.SignToken(thisUser.ID)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, models.TokenDTO{Token: tokenString})

}

func getToken(c *echo.Context) (string, error) {
	req := (*c).Request()
	authSlice := req.Header["Auth"]

	if len(authSlice) == 0 {
		return "", fmt.Errorf("token does not exist")
	}

	s := strings.Split(authSlice[0], " ")

	if len(s) != 2 || s[0] != "ut" {
		return "", fmt.Errorf("bad token format")
	}

	return s[1], nil
}

func generateMe(c *echo.Context, h *HandlerService) (*models.User, int, error) {
	tokenString, err := getToken(c)

	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("unathorized")
	}

	token, err := utils.VerifyToken(tokenString)

	if err != nil || !token.Valid {
		return nil, http.StatusUnauthorized, errors.New("unathorized")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["id"] == nil {
		return nil, http.StatusUnauthorized, errors.New("unathorized")
	}

	i := claims["id"].(string)
	userId, err := strconv.Atoi(i)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("this one's on us")
	}

	thisUser, err := FindSingleUser(h, userId)

	if err != nil || thisUser == nil {
		return nil, http.StatusUnauthorized, errors.New("unathorized")
	}

	if !thisUser.IsEmailVerified {
		return nil, http.StatusForbidden, errors.New("your email is not verified")
	}

	return thisUser, 0, nil

}

func (h *HandlerService) AuthGaurd(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		me, code, err := generateMe(&c, h)
		if err != nil {
			return WriteReponse(c, code, err.Error())
		}

		if me.Role != "user" {
			return WriteReponse(c, http.StatusForbidden, "forbidden resources")
		}

		cc := CustomContext{c, me}
		return next(cc)

	}

}

func (h *HandlerService) AdminGaurd(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		me, code, err := generateMe(&c, h)
		if err != nil || me == nil {
			return WriteReponse(c, code, err.Error())
		}

		if me.Role != "admin" {
			return WriteReponse(c, http.StatusForbidden, "forbidden resources")
		}

		cc := CustomContext{c, me}
		return next(cc)

	}
}

func GetMe(c *echo.Context) *models.User {
	return (*c).(CustomContext).User
}

func (h *HandlerService) HandleMe(c echo.Context) error {

	me := GetMe(&c)

	user := dataprocesslayer.ConvertToUserDto(me)

	return c.JSON(http.StatusOK, user)
}

func (h *HandlerService) HandleRequestChangePassword(c echo.Context) error {
	body := models.RequestChangePasswordDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "please provide a valid email")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "no user found")
	}

	if !thisUser.IsEmailVerified {
		return WriteReponse(c, http.StatusBadRequest, "this option is for those who have validated their emails")
	}

	code := CreateUUID()

	redisErr := h.redisStore.SetPasswordChangeRequest(thisUser.Email, code)
	if redisErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	go h.mailerService.SendPasswordChangeRequestEmail(thisUser, code)

	return WriteReponse(c, http.StatusOK, "check your email")

}

func (h *HandlerService) HandleVerifyChangePasswordRequest(c echo.Context) error {
	email := c.QueryParam("email")
	code := c.QueryParam("code")

	if email == "" || code == "" {
		return WriteReponse(c, http.StatusBadRequest, "insufficient data")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "user not found")
	}

	if val, err := h.redisStore.GetPasswordChangeRequest(thisUser.Email); err != nil || val != code {
		return WriteReponse(c, http.StatusBadRequest, "code doesnt match")
	}

	go h.redisStore.DeletePasswordChangeRequest(thisUser.Email)
	redisErr := h.redisStore.CreatePasswordChangePermission(thisUser.Email, code)
	if redisErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	responsetext := fmt.Sprintf("you can change your password at %v/auth/change-password", os.Getenv("DOMAIN"))
	return WriteReponse(c, http.StatusOK, responsetext)
}

func (h *HandlerService) HandleChangePassword(c echo.Context) error {
	body := models.ChangePasswordDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "please provide a valid email")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "no user found")
	}

	if val, err := h.redisStore.GetPasswordChangePermission(thisUser.Email); err != nil || val != body.Code {
		return WriteReponse(c, http.StatusForbidden, "access denied")
	}

	if err := h.store.UserRepo.UpdatePassword(body.Email, body.Password); err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	go h.redisStore.DelPasswordChangePermission(body.Email)

	return WriteReponse(c, http.StatusOK, "password changed successfully")
}

func (h *HandlerService) CreateAdmin(c echo.Context) error {

	body := models.CreateAdminDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "lack of data")
	}

	thisUser, _ := h.store.UserRepo.FindByEmail(body.Email)

	if thisUser != nil {
		return WriteReponse(c, http.StatusBadRequest, "this email already exists in the database")
	}

	var err error
	thisUser, err = h.store.UserRepo.Create(body.Name, body.Email, body.Password, "admin", true)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, err.Error())
	}

	user := dataprocesslayer.ConvertToUserDto(thisUser)
	return c.JSON(http.StatusOK, user)
}

func (h *HandlerService) GetSingleUser(c echo.Context) error {

	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	thisUser, err := FindSingleUser(h, userId)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}

	user := dataprocesslayer.ConvertToUserDto(thisUser)

	return c.JSON(http.StatusOK, user)
}

func (h *HandlerService) GetUsers(c echo.Context) error {

	text := c.QueryParam("text")
	p := c.QueryParam("page")
	l := c.QueryParam("limit")

	var page int
	var limit int
	var err error

	page, err = strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err = strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	users, fErr := h.store.UserRepo.FindUsers(text, page, limit)
	if fErr != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}

	var result []models.UserDto

	for _, s := range users {
		result = append(result, dataprocesslayer.ConvertToUserDto(s))
	}

	return c.JSON(http.StatusOK, result)
}
