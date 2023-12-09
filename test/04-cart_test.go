package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/thejixer/shop-api/internal/models"
)

/*

	1) u1 should be able to get its cart with 0 items in it
	2) should be able to add p1 to its card and
	have 1 as cartlength and readd p1 and still have the same result
	3) should be able to remove p1 and have 0 as cartlength
	4) should be able to add p1 and p2 and have 2 as cartlength
	5) should be able to checkout and have 0 as cartlength


*/

func getUserCart(token string, t *testing.T) *models.CartDto {

	res := Fetch("GET", "cart/", token, nil, t)
	cart := new(models.CartDto)
	json.Unmarshal(res, cart)

	return cart
}

func AddToCart(token string, productId, quantity int, t *testing.T) {
	AddToCartData := models.AddtoCartDto{
		ProductId: productId,
		Quantity:  quantity,
	}
	bodyString, err := json.Marshal(AddToCartData)
	if err != nil {
		t.Error("this one is on code")
	}
	body := []byte(bodyString)
	res := Fetch("POST", "cart/", token, bytes.NewBuffer(body), t)
	response := models.ResponseDTO{}
	json.Unmarshal(res, &response)
	expected := 201
	if response.StatusCode != expected {
		t.Errorf("expected to recieve status code %v but got %v", expected, response.StatusCode)
	}
}

func removeItemFromCart(token string, productId int, t *testing.T) {
	res := Fetch("DELETE", fmt.Sprintf("cart/%v", productId), token, nil, t)
	response := models.ResponseDTO{}
	json.Unmarshal(res, &response)
	expectedStatus := http.StatusAccepted
	if response.StatusCode != expectedStatus {
		t.Errorf("expected status of %v but got %v", expectedStatus, response.StatusCode)
	}
}

func chargeBalance(token string, amount float64, t *testing.T) {
	chargeBalanceData := models.ChargeBalanceDto{
		Amount: amount + 100,
	}
	bodyString, err := json.Marshal(chargeBalanceData)
	if err != nil {
		t.Error("this one is on code")
	}
	body := []byte(bodyString)

	res := Fetch("POST", "user/charge-balance", token, bytes.NewBuffer(body), t)
	response := models.ResponseDTO{}
	json.Unmarshal(res, &response)
	expectedStatus := http.StatusAccepted
	if response.StatusCode != expectedStatus {
		t.Errorf("expected status of %v but got %v", expectedStatus, response.StatusCode)
	}
}

func TestStartofCart(t *testing.T) {
	u1Token := GetContect("u1Token").(string)
	x := getUserCart(u1Token, t)
	expected := 0
	if len(x.Items) != expected {
		t.Errorf("expected to recieve %v items but got %v", expected, len(x.Items))
	}

}

func TestAddToCart(t *testing.T) {
	u1Token := GetContect("u1Token").(string)
	p1 := GetContect("p1").(models.ProductDto)

	AddToCart(u1Token, p1.Id, 1, t)

	myCart := getUserCart(u1Token, t)

	if len(myCart.Items) != 1 {
		t.Errorf("expected to recieve %v items but got %v", 1, len(myCart.Items))
	}

	AddToCart(u1Token, p1.Id, 3, t)

	myCart = getUserCart(u1Token, t)
	if len(myCart.Items) != 1 {
		t.Errorf("expected to recieve %v items but got %v", 1, len(myCart.Items))
	}

	if myCart.Items[0].Quantity != 3 {
		t.Errorf("expected to have %v of product with id of %v but got %v", 3, p1.Id, myCart.Items[0].Quantity)
	}
}

func TestRemoveFromCart(t *testing.T) {
	u1Token := GetContect("u1Token").(string)
	p1 := GetContect("p1").(models.ProductDto)

	removeItemFromCart(u1Token, p1.Id, t)

	myCart := getUserCart(u1Token, t)

	expetedItemCount := 0
	if len(myCart.Items) != expetedItemCount {
		t.Errorf("expected to recieve %v items but got %v", expetedItemCount, len(myCart.Items))
	}

}

func TestAddToCartAgain(t *testing.T) {

	u1Token := GetContect("u1Token").(string)
	p1 := GetContect("p1").(models.ProductDto)
	p2 := GetContect("p2").(models.ProductDto)

	AddToCart(u1Token, p1.Id, 1, t)
	AddToCart(u1Token, p2.Id, 1, t)

	myCart := getUserCart(u1Token, t)

	expetedItemCount := 2
	if len(myCart.Items) != expetedItemCount {
		t.Errorf("expected to recieve %v items but got %v", expetedItemCount, len(myCart.Items))
	}

}

func TestCheckout(t *testing.T) {
	u1Token := GetContect("u1Token").(string)
	a1 := GetContect("a1").(*models.AddressDto)
	myCart := getUserCart(u1Token, t)

	CheckOutData := models.CheckOutDto{
		AddressId: a1.Id,
	}
	bodyString, err := json.Marshal(CheckOutData)
	if err != nil {
		t.Error("this one is on code")
	}
	body := []byte(bodyString)
	res := Fetch("POST", "checkout", u1Token, bytes.NewBuffer(body), t)
	response := models.ResponseDTO{}
	json.Unmarshal(res, &response)
	expected := http.StatusBadRequest
	if response.StatusCode != expected {
		t.Errorf("expected to recieve status code %v but got %v", expected, response.StatusCode)
	}

	chargeBalance(u1Token, myCart.TotalPrice, t)

	res = Fetch("POST", "checkout", u1Token, bytes.NewBuffer(body), t)
	response = models.ResponseDTO{}
	json.Unmarshal(res, &response)
	expected = http.StatusCreated
	if response.StatusCode != expected {
		t.Errorf("expected to recieve status code %v but got %v", expected, response.StatusCode)
	}

	myCart = getUserCart(u1Token, t)

	expetedItemCount := 0
	if len(myCart.Items) != expetedItemCount {
		t.Errorf("expected to recieve %v items but got %v", expetedItemCount, len(myCart.Items))
	}

}
