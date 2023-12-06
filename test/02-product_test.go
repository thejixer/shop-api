package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/thejixer/shop-api/internal/models"
)

/*

	1) create a product
	2) edit that product
	3) get that single product
	4) get all products

*/

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
	res = Fetch("POST", "product/", adminToken, bytes.NewBuffer(body), t)
	product2 := models.ProductDto{}
	json.Unmarshal(res, &product2)
	SetContext("p2", product2)
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
	result := models.LL_ProductDto{}
	json.Unmarshal(res, &result)
	if result.Total < 1 {
		t.Errorf("expected to get at least 1 product back but got nothing")
	}
}
