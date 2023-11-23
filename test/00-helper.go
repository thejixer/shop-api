package test

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
)

type User struct {
	Name     string
	Email    string
	Password string
	Token    string
}

var Domain string
var customContext = map[string]any{}

func SetContext(key string, value any) {
	customContext[key] = value
}

func GetContect(key string) any {
	return customContext[key]
}

func GenerateRandomSingleDigit() int {
	return rand.Intn(10)
}

func GenerateNumericString(n int) string {
	res := ""
	for i := n; i < n; i++ {
		res = fmt.Sprintf("%v%v", res, GenerateRandomSingleDigit())
	}
	return res
}

func Fetch(
	method,
	url,
	token string,
	body io.Reader,
	t *testing.T,
) []byte {

	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("%v/%v", Domain, url), body)
	if err != nil {
		t.Error("server aint reachable")
	}
	if token != "" {
		authHeader := fmt.Sprintf("ut %v", token)
		req.Header.Add("auth", authHeader)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Error("server aint reachable")
	}

	data, er := io.ReadAll(resp.Body)
	if er != nil {
		t.Error("server aint reachable")
	}
	return data
}
