package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/thejixer/shop-api/internal/models"
)

func TestAdminLogin(t *testing.T) {

	bodyString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, os.Getenv("MAIN_ADMIN_EMAIL"), os.Getenv("MAIN_ADMIN_PASSWORD"))

	body := []byte(bodyString)
	res := Fetch("POST", "auth/login", "", bytes.NewBuffer(body), t)

	result := models.TokenDTO{}
	json.Unmarshal(res, &result)
	if len(result.Token) < 10 {
		t.Errorf("expected to get a token but got wierd stuff")
	}

	SetContext("admin_token", result.Token)

}

func TestCreateProduct(t *testing.T) {
	adminToken := fmt.Sprintf("%v", GetContect("admin_token"))

	bodyString := fmt.Sprintf(`{ 
    "title": "Test Product",
    "price": 99.99,
    "description": "Just a product for Test",
    "quantity": 999
	}`)
	body := []byte(bodyString)
	res := Fetch("POST", "product/", "", bytes.NewBuffer(body), t)
	result := models.ResponseDTO{}
	json.Unmarshal(res, &result)

	expected := 401
	if result.StatusCode != expected {
		t.Errorf("expected to get status code of %v but got %v", expected, result.StatusCode)
	}

	res = Fetch("POST", "product/", adminToken, bytes.NewBuffer(body), t)
	product := models.ProductDto{}
	json.Unmarshal(res, &product)
	expectedTitle := "Test Product"
	if product.Title != expectedTitle {
		t.Errorf("expected to get status code of %v but got %v", expectedTitle, product.Title)
	}

	SetContext("p1", product)

}

func TestEditProduct(t *testing.T) {
	adminToken := fmt.Sprintf("%v", GetContect("admin_token"))
	p1 := GetContect("p1").(models.ProductDto)
	expectedPrice := 109.99

	bodyString := fmt.Sprintf(`{ 
    "title": "[Edited] Test Product",
    "price": %v,
    "description": "[Edited] Just a product for Test",
    "quantity": 999
	}`, expectedPrice)
	body := []byte(bodyString)
	res := Fetch("POST", fmt.Sprintf("product/%v", p1.Id), adminToken, bytes.NewBuffer(body), t)
	result := models.ProductDto{}
	json.Unmarshal(res, &result)

	t.Log(result)

	if result.Price != expectedPrice {
		t.Errorf("expected to get the price of %v but got %v", expectedPrice, result.Price)
	}
	SetContext("p1", result)

}

func TestGetSingleProduct(t *testing.T) {
	p1 := GetContect("p1").(models.ProductDto)
	res := Fetch("GET", fmt.Sprintf("product/%v", p1.Id), "", nil, t)
	result := models.ProductDto{}
	json.Unmarshal(res, &result)

	if result.Title != p1.Title {
		t.Errorf("something's wrong")
	}

}

func TestGetProducts(t *testing.T) {
	res := Fetch("GET", "product/?", "", nil, t)
	result := []models.ProductDto{}
	json.Unmarshal(res, &result)
	if len(result) < 1 {
		t.Errorf("expected to get at least 1 product back but got nothing")
	}
}
