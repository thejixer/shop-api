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

	1) u1 should be able to get his addresses and the total should be 0
	2) should be able to create a new address and the address count should be 1
	3) should be able to delete the last address and the addresscount should be 0
	4) he will create another address just because we need an address later

*/

func getAddressess(token string, t *testing.T) []*models.AddressDto {

	res := Fetch("GET", "address/", token, nil, t)
	addresses := []*models.AddressDto{}
	json.Unmarshal(res, &addresses)

	return addresses
}

func createAddress(token string, t *testing.T) {
	createAddressData := models.CreateAddressDto{
		Title:         "test address",
		Lon:           33,
		Lat:           33,
		Address:       "the obvious door",
		RecieverName:  "me",
		RecieverPhone: "09123344555",
	}
	bodyString, err := json.Marshal(createAddressData)
	if err != nil {
		t.Error("this one is on code")
	}
	body := []byte(bodyString)
	res := Fetch("POST", "address/", token, bytes.NewBuffer(body), t)
	response := models.AddressDto{}
	json.Unmarshal(res, &response)
	expectedTitle := createAddressData.Title
	if response.Title != expectedTitle {
		t.Errorf("expected to recieve title of %v but got %v", expectedTitle, response.Title)
	}
}

func TestAddressesStart(t *testing.T) {
	u1Token := GetContect("u1Token").(string)
	myAddresses := getAddressess(u1Token, t)
	expectedLength := 0
	if len(myAddresses) != expectedLength {
		t.Errorf("expected to recieve %v addresses but got %v", expectedLength, myAddresses)
	}
}

func TestCreateAddress(t *testing.T) {
	u1Token := GetContect("u1Token").(string)
	createAddress(u1Token, t)
	createAddress(u1Token, t)
	myAddresses := getAddressess(u1Token, t)
	expectedLength := 2
	if len(myAddresses) != expectedLength {
		t.Errorf("expected to recieve %v addresses but got %v", expectedLength, myAddresses)
	}
	SetContext("a1", myAddresses[0])
	SetContext("a2", myAddresses[1])
}

func TestDeleteAddress(t *testing.T) {
	u1Token := GetContect("u1Token").(string)
	a2 := GetContect("a2").(*models.AddressDto)
	res := Fetch("DELETE", fmt.Sprintf("address/%v", a2.Id), u1Token, nil, t)
	response := models.ResponseDTO{}
	json.Unmarshal(res, &response)

	expectedStatus := http.StatusAccepted
	if response.StatusCode != expectedStatus {
		t.Errorf("expected to get status code of %v but got %v", expectedStatus, response.StatusCode)
	}

	myAddresses := getAddressess(u1Token, t)
	expectedLength := 1
	if len(myAddresses) != expectedLength {
		t.Errorf("expected to recieve %v addresses but got %v", expectedLength, len(myAddresses))
	}

}
