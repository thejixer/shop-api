package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/thejixer/shop-api/models"
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

type User struct {
	Name     string
	Email    string
	Password string
	Token    string
}

var domain string
var customContext = map[string]any{}

func setContext(key string, value any) {
	customContext[key] = value
}

func getContect(key string) any {
	return customContext[key]
}

func generateRandomSingleDigit() int {
	return rand.Intn(10)
}

func generateNumericString(n int) string {
	res := ""
	for i := n; i < n; i++ {
		res = fmt.Sprintf("%v%v", res, generateRandomSingleDigit())
	}
	return res
}

func Fetch(url, method string, body io.Reader) ([]byte, error) {

	switch method {
	case "GET":
		resp, err := http.Get(url)
		if err != nil {
			return nil, errors.New("server unreachable")
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New("server unreachable")
		}

		return data, nil
	case "POST":
		resp, err := http.Post(url, "application/json", body)
		if err != nil {
			return nil, errors.New("server unreachable")
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New("server unreachable")
		}

		return data, nil
	default:
		return nil, errors.New("bad method")
	}

}

func init() {
	godotenv.Load("../.env")
	domain = os.Getenv("DOMAIN")
}

func TestRoot(t *testing.T) {

	body, err := Fetch(domain, "GET", nil)

	if err != nil {
		t.Error("server unreachable")
	}

	expected := "Hello, World from shop-api!"
	if x := string(body); x != expected {
		t.Errorf("expected %s but got %s", expected, x)
	}

}

func TestSignUp(t *testing.T) {

	email := fmt.Sprintf("test-%v%v@test.com", time.Now().UnixMilli(), generateNumericString(4))

	u1 := User{}
	u1.Name = "u1 user"
	u1.Email = email
	u1.Password = "123456"

	setContext("u1", u1)

	bodyString := fmt.Sprintf(`{"name": "%v", "email": "%v", "password": "%v" }`, u1.Name, u1.Email, u1.Password)

	body := []byte(bodyString)

	res, err := Fetch(domain+"/auth/signup", "POST", bytes.NewBuffer(body))
	if err != nil {
		t.Error("server aint reachable")
	}

	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)

	expected := 200
	if result.StatusCode != expected {
		t.Errorf("expected %v but got %v", expected, result.StatusCode)
	}

}

func TestFailedLogin(t *testing.T) {
	u1 := getContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, u1.Email, u1.Password)

	body := []byte(bodyString)
	res, err := Fetch(domain+"/auth/login", "POST", bytes.NewBuffer(body))
	if err != nil {
		t.Error("server aint reachable")
	}
	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)
	expected := 400
	if result.StatusCode != expected {
		t.Errorf("expected %v but got %v", expected, result.StatusCode)
	}
}

func TestVerifyEmail(t *testing.T) {

	u1 := getContect("u1").(User)
	url := fmt.Sprintf("%v/auth/verify-email?email=%v&code=%v", domain, u1.Email, "1111")
	res, err := Fetch(url, "GET", nil)
	if err != nil {
		t.Error("server aint reachable")
	}
	data := models.TokenDTO{}
	json.Unmarshal(res, &data)

	if len(data.Token) < 10 {
		t.Errorf("expected to get a token but got %s", data.Token)
	}

	u1.Token = data.Token
	setContext("u1", u1)

}

func TestLogin(t *testing.T) {
	u1 := getContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, u1.Email, u1.Password)

	body := []byte(bodyString)
	res, err := Fetch(domain+"/auth/login", "POST", bytes.NewBuffer(body))
	if err != nil {
		t.Error("server aint reachable")
	}
	result := models.TokenDTO{}
	json.Unmarshal(res, &result)
	if len(result.Token) < 10 {
		t.Errorf("expected to get a token but got %v", result.Token)
	}

}

func TestMe(t *testing.T) {
	u1 := getContect("u1").(User)

	client := &http.Client{}

	req, err := http.NewRequest("POST", domain+"/auth/me", nil)
	if err != nil {
		t.Error("server aint reachable")
	}
	authHeader := fmt.Sprintf("ut %v", u1.Token)
	req.Header.Add("auth", authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Error("server aint reachable")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error("server aint reachable")
	}

	me := models.UserDto{}
	json.Unmarshal(body, &me)

	if me.Email != u1.Email {
		t.Error("server aint reachable")
	}

}

func TestRequestChangePassword(t *testing.T) {
	u1 := getContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v" }`, u1.Email)
	body := []byte(bodyString)
	res, err := Fetch(domain+"/auth/request-change-password", "POST", bytes.NewBuffer(body))
	if err != nil {
		t.Error("server aint reachable")
	}
	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)

	expected := 200
	if result.StatusCode != expected {
		t.Errorf("result.msg : %v", result.Msg)
		t.Errorf("expected to get status code %v but got %v", expected, result.StatusCode)
	}

	url := fmt.Sprintf("%v/auth/verify-changepassword-request?email=%v&code=%v", domain, u1.Email, "1111")
	res, err = Fetch(url, "GET", nil)
	if err != nil {
		t.Error("server aint reachable")
	}
	result = models.ResponseDTO{}
	json.Unmarshal(res, &result)

	if result.StatusCode != expected {
		t.Errorf("result.msg : %v", result.Msg)
		t.Errorf("expected to get status code %v but got %v", expected, result.StatusCode)
	}

	u1.Password = "123456789"
	setContext("u1", u1)
	bodyString = fmt.Sprintf(`{ "email": "%v", "password": "%v", "code": "%v" }`, u1.Email, u1.Password, "1111")
	body = []byte(bodyString)
	res, err = Fetch(domain+"/auth/change-password", "POST", bytes.NewBuffer(body))
	if err != nil {
		t.Error("server aint reachable")
	}
	result = models.ResponseDTO{}
	json.Unmarshal(res, &result)

	if result.StatusCode != expected {
		t.Errorf("expected to get status code %v but got %v", expected, result.StatusCode)
	}

}

func TestFailedLoginWithOldPassword(t *testing.T) {
	u1 := getContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, u1.Email, "123456")

	body := []byte(bodyString)
	res, err := Fetch(domain+"/auth/login", "POST", bytes.NewBuffer(body))
	if err != nil {
		t.Error("server aint reachable")
	}
	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)
	expected := 400
	if result.StatusCode != expected {
		t.Errorf("expected %v but got %v", expected, result.StatusCode)
	}
}

func TestSuccessfulLoginwithNewPassword(t *testing.T) {
	u1 := getContect("u1").(User)
	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, u1.Email, u1.Password)

	body := []byte(bodyString)
	res, err := Fetch(domain+"/auth/login", "POST", bytes.NewBuffer(body))
	if err != nil {
		t.Error("server aint reachable")
	}
	result := models.TokenDTO{}
	json.Unmarshal(res, &result)
	if len(result.Token) < 10 {
		t.Errorf("expected to get a token but got wierd stuff")
	}
}
