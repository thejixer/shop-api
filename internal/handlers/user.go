package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	thisUser, err = h.store.UserRepo.Create(body.Name, body.Email, body.Password, "user", false, []string{})
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
		fmt.Println(err)
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
		fmt.Println(err)
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

func GetMe(c *echo.Context) (*models.User, error) {
	me := (*c).(CustomContext).User
	if me == nil {
		return nil, errors.New("unathorized")
	}
	return me, nil
}

func (h *HandlerService) HandleMe(c echo.Context) error {

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	user := dataprocesslayer.ConvertToUserDto(me)

	return c.JSON(http.StatusOK, user)
}

func (h *HandlerService) HandleAdminMe(c echo.Context) error {
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	admin := dataprocesslayer.ConvertToAdminDto(me)

	return c.JSON(http.StatusOK, admin)
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

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	body := models.CreateAdminDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid data")
	}

	if hasPermission := PermissionChecker(me, "master"); !hasPermission {
		return WriteReponse(c, http.StatusForbidden, "forbidden resources")
	}

	thisUser, err := h.store.UserRepo.Create(body.Name, body.Email, body.Password, "admin", true, body.Permissions)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return WriteReponse(c, http.StatusBadRequest, "email already exists")
		}
		return WriteReponse(c, http.StatusInternalServerError, "oops, this one's on us")
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

	users, count, fErr := h.store.UserRepo.FindUsers(text, page, limit)
	if fErr != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}

	result := dataprocesslayer.ConvertToLLUserDto(users, count)

	return c.JSON(http.StatusOK, result)
}

func (h *HandlerService) ChargeBalance(c echo.Context) error {
	// for simplicity reasons we dont have payment systems and users can charge their balance as they wish
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	body := models.ChargeBalanceDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil || body.Amount <= 0 {
		return WriteReponse(c, http.StatusBadRequest, "invalid input")
	}

	if err := h.store.UserRepo.ChargeBalance(me.ID, body.Amount); err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "oops, this one's on us")
	}

	return WriteReponse(c, http.StatusAccepted, "successfully charged your account")

}

func (h *HandlerService) UpdatePermissions(c echo.Context) error {
	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	body := models.UpdatePermissionDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "invalid input")
	}

	if hasPermission := PermissionChecker(me, "master"); !hasPermission {
		return WriteReponse(c, http.StatusForbidden, "forbidden resources")
	}

	if thatUser, err := h.store.UserRepo.FindById(body.UserId); err != nil || thatUser.Role != "admin" {
		return WriteReponse(c, http.StatusBadRequest, "bad request")
	}

	if err := h.store.UserRepo.UpdatePermissions(body.UserId, body.Permissions); err != nil {
		fmt.Println(err)
		return WriteReponse(c, http.StatusInternalServerError, "oops, this one's on us")
	}

	go h.redisStore.DelUser(body.UserId)

	return WriteReponse(c, http.StatusAccepted, "successfully updated users permissions")

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
