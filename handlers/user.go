package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/shop-api/data-process-layer"
	"github.com/thejixer/shop-api/models"
	"github.com/thejixer/shop-api/utils"
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
		return WriteReponse(&c, http.StatusBadRequest, "lack of data")
	}

	thisUser, _ := h.store.UserRepo.FindByEmail(body.Email)

	if thisUser != nil {
		return WriteReponse(&c, http.StatusBadRequest, "this email already exists in the database")
	}

	var err error
	thisUser, err = h.store.UserRepo.Create(body)
	if err != nil {
		return WriteReponse(&c, http.StatusBadRequest, err.Error())
	}

	verificationCode := CreateUUID()

	redisErr := h.redisStore.SetEmailVerificationCode(thisUser.Email, verificationCode)
	if redisErr != nil {
		return WriteReponse(&c, http.StatusInternalServerError, "this is on us, please try again")
	}

	go h.mailerService.SendVerificationEmail(thisUser, verificationCode)

	return WriteReponse(&c, http.StatusOK, "please check your email to verify your email")
}

func (h *HandlerService) HandleRequestVerificationEmail(c echo.Context) error {
	body := models.RequestVerificationEmailDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(&c, http.StatusBadRequest, "lack of data")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(&c, http.StatusNotFound, "no user found")
	}

	if thisUser.IsEmailVerified {
		return WriteReponse(&c, http.StatusBadRequest, "your email has already been verified")
	}

	verificationCode := CreateUUID()
	redisErr := h.redisStore.SetEmailVerificationCode(thisUser.Email, verificationCode)
	if redisErr != nil {
		return WriteReponse(&c, http.StatusInternalServerError, "this is on us, please try again")
	}

	go h.mailerService.SendVerificationEmail(thisUser, verificationCode)

	return WriteReponse(&c, http.StatusOK, "please check your email to verify your email")
}

func (h *HandlerService) HandleEmailVerification(c echo.Context) error {
	email := c.QueryParam("email")
	verificationCode := c.QueryParam("code")

	fmt.Println("*****************")
	fmt.Println("email : ", email)

	if email == "" || verificationCode == "" {
		return WriteReponse(&c, http.StatusBadRequest, "insufficient data")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(email)

	if err != nil {
		return WriteReponse(&c, http.StatusNotFound, "user not found")
	}

	if val, err := h.redisStore.GetEmailVerificationCode(thisUser.Email); err != nil || val != verificationCode {
		return WriteReponse(&c, http.StatusBadRequest, "code doesnt match")
	}

	updateErr := h.store.UserRepo.VerifyEmail(email)
	go h.redisStore.DeleteEmailVerificationCode(email)
	if updateErr != nil {
		return WriteReponse(&c, http.StatusInternalServerError, "this one's on us")
	}

	tokenString, err := utils.SignToken(thisUser.ID)
	if err != nil {
		return WriteReponse(&c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, models.TokenDTO{Token: tokenString})
}

func (h *HandlerService) HandleLogin(c echo.Context) error {
	body := models.LoginDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(&c, http.StatusBadRequest, "lack of data")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(body.Email)
	if err != nil {
		return WriteReponse(&c, http.StatusNotFound, "no user found")
	}

	if !thisUser.IsEmailVerified {
		return WriteReponse(&c, http.StatusBadRequest, "your email is not verified")
	}

	if match := utils.CheckPasswordHash(body.Password, thisUser.Password); !match {
		return WriteReponse(&c, http.StatusBadRequest, "password doesnt match")
	}

	tokenString, err := utils.SignToken(thisUser.ID)
	if err != nil {
		return WriteReponse(&c, http.StatusInternalServerError, err.Error())
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

func (h *HandlerService) AuthGaurd(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		tokenString, err := getToken(&c)

		if err != nil {
			return WriteReponse(&c, http.StatusUnauthorized, "unathorized")
		}

		token, err := utils.VerifyToken(tokenString)

		if err != nil || !token.Valid {
			return WriteReponse(&c, http.StatusUnauthorized, "unathorized")
		}

		claims := token.Claims.(jwt.MapClaims)

		if claims["id"] == nil {
			return WriteReponse(&c, http.StatusUnauthorized, "unathorized")
		}

		i := claims["id"].(string)
		userId, err := strconv.Atoi(i)

		if err != nil {
			return WriteReponse(&c, http.StatusInternalServerError, "this one's on us")
		}

		thisUser, err := h.store.UserRepo.FindById(userId)

		if err != nil || thisUser == nil {
			return WriteReponse(&c, http.StatusUnauthorized, "unathorized")
		}

		if !thisUser.IsEmailVerified {
			return WriteReponse(&c, http.StatusBadRequest, "your email is not verified")
		}

		cc := CustomContext{c, thisUser}
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
		return WriteReponse(&c, http.StatusBadRequest, "please provide a valid email")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(&c, http.StatusNotFound, "no user found")
	}

	if !thisUser.IsEmailVerified {
		return WriteReponse(&c, http.StatusBadRequest, "this option is for those who have validated their emails")
	}

	code := CreateUUID()

	redisErr := h.redisStore.SetPasswordChangeRequest(thisUser.Email, code)
	if redisErr != nil {
		return WriteReponse(&c, http.StatusInternalServerError, "this is on us, please try again")
	}

	go h.mailerService.SendPasswordChangeRequestEmail(thisUser, code)

	return WriteReponse(&c, http.StatusOK, "check your email")

}

func (h *HandlerService) HandleVerifyChangePasswordRequest(c echo.Context) error {
	email := c.QueryParam("email")
	code := c.QueryParam("code")

	if email == "" || code == "" {
		return WriteReponse(&c, http.StatusBadRequest, "insufficient data")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(email)

	if err != nil {
		return WriteReponse(&c, http.StatusNotFound, "user not found")
	}

	if val, err := h.redisStore.GetPasswordChangeRequest(thisUser.Email); err != nil || val != code {
		return WriteReponse(&c, http.StatusBadRequest, "code doesnt match")
	}

	go h.redisStore.DeletePasswordChangeRequest(thisUser.Email)
	redisErr := h.redisStore.CreatePasswordChangePermission(thisUser.Email, code)
	if redisErr != nil {
		return WriteReponse(&c, http.StatusInternalServerError, "this is on us, please try again")
	}

	responsetext := fmt.Sprintf("you can change your password at %v/auth/change-password", os.Getenv("DOMAIN"))
	return WriteReponse(&c, http.StatusOK, responsetext)

}

func (h *HandlerService) HandleChangePassword(c echo.Context) error {
	body := models.ChangePasswordDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(&c, http.StatusBadRequest, "please provide a valid email")
	}

	thisUser, err := h.store.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(&c, http.StatusNotFound, "no user found")
	}

	if val, err := h.redisStore.GetPasswordChangePermission(thisUser.Email); err != nil || val != body.Code {
		return WriteReponse(&c, http.StatusForbidden, "access denied")
	}

	if err := h.store.UserRepo.UpdatePassword(body.Email, body.Password); err != nil {
		return WriteReponse(&c, http.StatusInternalServerError, "this one's on us")
	}

	go h.redisStore.DelPasswordChangePermission(body.Email)

	return WriteReponse(&c, http.StatusOK, "password changed successfully")

}
