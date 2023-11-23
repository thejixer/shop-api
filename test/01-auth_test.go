package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/thejixer/shop-api/internal/models"
)

/*

	Authentication Scneario Tests:
	1) Signup a new user u1
	2) they should not be able to login since they haven't verified their email
	3) u1 verifies their email
	4) u1 should be able to login and get a token
	5) /auth/me should return u1's personal data with the token
	6) u1 can change their password
	7) they obviously should not be able to login with their old password
	8) they can login with the new password and get their new token

*/

func init() {
	godotenv.Load("../.env")
	Domain = os.Getenv("DOMAIN")
}

func TestRoot(t *testing.T) {

	res := Fetch("GET", "", "", nil, t)

	expected := "Hello, World from shop-api!"
	if x := string(res); x != expected {
		t.Errorf("expected %s but got %s", expected, x)
	}

}

func TestSignUp(t *testing.T) {

	email := fmt.Sprintf("test-%v%v@test.com", time.Now().UnixMilli(), GenerateNumericString(4))

	u1 := User{}
	u1.Name = "u1 user"
	u1.Email = email
	u1.Password = "123456"

	SetContext("u1", u1)

	bodyString := fmt.Sprintf(`{"name": "%v", "email": "%v", "password": "%v" }`, u1.Name, u1.Email, u1.Password)

	body := []byte(bodyString)

	res := Fetch("POST", "auth/signup", "", bytes.NewBuffer(body), t)

	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)

	expected := 200
	if result.StatusCode != expected {
		t.Errorf("expected %v but got %v", expected, result.StatusCode)
	}

}

func TestFailedLogin(t *testing.T) {
	u1 := GetContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, u1.Email, u1.Password)

	body := []byte(bodyString)
	res := Fetch("POST", "auth/login", "", bytes.NewBuffer(body), t)

	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)
	expected := 401
	if result.StatusCode != expected {
		t.Errorf("expected %v but got %v", expected, result.StatusCode)
	}
}

func TestVerifyEmail(t *testing.T) {

	u1 := GetContect("u1").(User)
	url := fmt.Sprintf("auth/verify-email?email=%v&code=%v", u1.Email, "1111")
	res := Fetch("GET", url, "", nil, t)

	data := models.TokenDTO{}
	json.Unmarshal(res, &data)

	if len(data.Token) < 10 {
		t.Errorf("expected to get a token but got %s", data.Token)
	}

	u1.Token = data.Token
	SetContext("u1", u1)

}

func TestLogin(t *testing.T) {
	u1 := GetContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, u1.Email, u1.Password)

	body := []byte(bodyString)
	res := Fetch("POST", "auth/login", "", bytes.NewBuffer(body), t)

	result := models.TokenDTO{}
	json.Unmarshal(res, &result)
	if len(result.Token) < 10 {
		t.Errorf("expected to get a token but got %v", result.Token)
	}

}

func TestMe(t *testing.T) {
	u1 := GetContect("u1").(User)

	res := Fetch("POST", "auth/me", u1.Token, nil, t)

	me := models.UserDto{}
	json.Unmarshal(res, &me)

	if me.Email != u1.Email {
		t.Error("expected to get u1's details but didn't")
	}

}

func TestRequestChangePassword(t *testing.T) {
	u1 := GetContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v" }`, u1.Email)
	body := []byte(bodyString)
	res := Fetch("POST", "auth/request-change-password", "", bytes.NewBuffer(body), t)

	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)

	expected := 200
	if result.StatusCode != expected {
		t.Errorf("result.msg : %v", result.Msg)
		t.Errorf("expected to get status code %v but got %v", expected, result.StatusCode)
	}

	url := fmt.Sprintf("auth/verify-changepassword-request?email=%v&code=%v", u1.Email, "1111")
	res = Fetch("GET", url, "", nil, t)

	result = models.ResponseDTO{}
	json.Unmarshal(res, &result)

	if result.StatusCode != expected {
		t.Errorf("result.msg : %v", result.Msg)
		t.Errorf("expected to get status code %v but got %v", expected, result.StatusCode)
	}

	u1.Password = "123456789"
	SetContext("u1", u1)
	bodyString = fmt.Sprintf(`{ "email": "%v", "password": "%v", "code": "%v" }`, u1.Email, u1.Password, "1111")
	body = []byte(bodyString)
	res = Fetch("POST", "auth/change-password", "", bytes.NewBuffer(body), t)

	result = models.ResponseDTO{}
	json.Unmarshal(res, &result)

	if result.StatusCode != expected {
		t.Errorf("expected to get status code %v but got %v", expected, result.StatusCode)
	}

}

func TestFailedLoginWithOldPassword(t *testing.T) {
	u1 := GetContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, u1.Email, "123456")

	body := []byte(bodyString)
	res := Fetch("POST", "auth/login", "", bytes.NewBuffer(body), t)

	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)
	expected := 400
	if result.StatusCode != expected {
		t.Errorf("expected %v but got %v", expected, result.StatusCode)
	}
}

func TestSuccessfulLoginwithNewPassword(t *testing.T) {
	u1 := GetContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, u1.Email, u1.Password)

	body := []byte(bodyString)
	res := Fetch("POST", "auth/login", "", bytes.NewBuffer(body), t)

	result := models.TokenDTO{}
	json.Unmarshal(res, &result)
	if len(result.Token) < 10 {
		t.Errorf("expected to get a token but got wierd stuff")
	}
}
