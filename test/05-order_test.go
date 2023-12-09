package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/thejixer/shop-api/internal/models"
)

/*

	0) 3 new admins should be created with different permissions

	1) u1 can get his orders with my-orders api and expect to have 1 order in it
	2) u1 and the admin should be able to get single order of O1
	3) admin should be able to see O1 in the get orders query by admin
	4) u-office should be able to verify order O1
	5) u-stock should be able to package the order O1
	6) u-office should be able to send the order O1
	7) u-shipper should be able to deliver the order O1 providing the code
	that the user gives him

*/

func createAdminAndReturnToken(masterToken string, Permissions []string, t *testing.T) (string, error) {
	newAdminData := models.CreateAdminDTO{
		Name:        "admin test",
		Email:       fmt.Sprintf("test-%v%v@test.com", time.Now().UnixMilli(), GenerateNumericString(4)),
		Password:    "123456",
		Permissions: Permissions,
	}
	bodyString, err := json.Marshal(newAdminData)
	if err != nil {
		t.Error("this one is on code")
		return "", errors.New("this one is on code")
	}
	body := []byte(bodyString)
	res := Fetch("POST", "admin/create", masterToken, bytes.NewBuffer(body), t)
	newAdmin := models.AdminDto{}
	json.Unmarshal(res, &newAdmin)

	loginString := fmt.Sprintf(`{ "email": "%v", "password": "%v" }`, newAdminData.Email, newAdminData.Password)
	body = []byte(loginString)
	res = Fetch("POST", "auth/login", "", bytes.NewBuffer(body), t)
	result := models.TokenDTO{}
	json.Unmarshal(res, &result)
	if len(result.Token) < 10 {
		t.Errorf("expected to get a token but got wierd stuff")
		return "", errors.New("expected to get a token but got wierd stuff")
	}

	return result.Token, nil

}

func getMyOrders(token string, t *testing.T) (*models.LL_OrderDto, error) {
	res := Fetch("GET", "order/my-orders?page=0&limit=10", token, nil, t)
	result := new(models.LL_OrderDto)
	json.Unmarshal(res, result)

	if strings.Contains(string(res), "unathorized") {
		return nil, errors.New("unathorized")
	}
	if strings.Contains(string(res), "forbidden") {
		return nil, errors.New("forbiden")
	}

	return result, nil
}

func getSingleOrder(token string, id int, t *testing.T) (*models.OrderDto, error) {
	res := Fetch("GET", fmt.Sprintf("order/single/%v", id), token, nil, t)
	thisOrder := new(models.OrderDto)
	json.Unmarshal(res, thisOrder)
	if strings.Contains(string(res), "forbidden") {
		return nil, errors.New("forbiden")
	}
	if strings.Contains(string(res), "unathorized") {
		return nil, errors.New("unathorized")
	}
	if thisOrder.Id == 0 {
		return nil, errors.New("bad request")
	}
	return thisOrder, nil
}

func adminGetOrders(token string, page, limit int, status string, userId int, t *testing.T) (*models.LL_OrderDto, error) {
	res := Fetch("GET", fmt.Sprintf("order/?page=%v&limit=%v&status=%v&userId=%v", page, limit, status, userId), token, nil, t)
	result := new(models.LL_OrderDto)
	json.Unmarshal(res, result)
	if strings.Contains(string(res), "unathorized") {
		return nil, errors.New("unathorized")
	}
	if strings.Contains(string(res), "forbidden") {
		return nil, errors.New("forbiden")
	}
	if result.Total == 0 {
		return nil, errors.New("bad request")
	}
	return result, nil
}

func TestPrepareOrderTests(t *testing.T) {
	adminToken := GetContect("admin_token").(string)

	officeToken, err := createAdminAndReturnToken(adminToken, []string{"backoffice"}, t)
	if err != nil {
		t.Errorf(err.Error())
	}
	SetContext("officeToken", officeToken)

	stockToken, err := createAdminAndReturnToken(adminToken, []string{"stock"}, t)
	if err != nil {
		t.Errorf(err.Error())
	}
	SetContext("stockToken", stockToken)

	shiperToken, err := createAdminAndReturnToken(adminToken, []string{"shipper"}, t)
	if err != nil {
		t.Errorf(err.Error())
	}
	SetContext("shiperToken", shiperToken)
}

func TestMyOrders(t *testing.T) {
	u1Token := GetContect("u1Token").(string)

	data, err := getMyOrders("", t)
	if err == nil {
		t.Error("expected to return an error but got no error")
	}

	data, err = getMyOrders(u1Token, t)

	if data.Total != 1 {
		t.Errorf("expected %v total orders but got %v", 1, data.Total)
	}

	expectedStatus := "created"
	if data.Result[0].Status != expectedStatus {
		t.Errorf("expected status of %v but got this wierd : %v", expectedStatus, data.Result[0].Status)
	}

	SetContext("o1", data.Result[0])
}

func TestSingleOder(t *testing.T) {
	u1Token := GetContect("u1Token").(string)
	adminToken := GetContect("admin_token").(string)
	o1 := GetContect("o1").(models.OrderDto)

	o1Data, err := getSingleOrder("", o1.Id, t)
	if err == nil {
		t.Error("expected to get an error ")
	}

	o1Data, err = getSingleOrder(u1Token, o1.Id, t)
	if err != nil || o1Data.Id != o1.Id {
		t.Errorf("expected to get order but didn't")
	}

	o1Data, err = getSingleOrder(adminToken, o1.Id, t)
	if err != nil || o1Data.Id != o1.Id {
		t.Errorf("expected to get order but didn't")
	}

}

func TestAdminGetOrders(t *testing.T) {
	adminToken := GetContect("admin_token").(string)
	u1Data := GetContect("u1Data").(models.UserDto)
	data, err := adminGetOrders("", 0, 10, "created", u1Data.ID, t)
	if err == nil {
		t.Errorf("expected to return an error but got no error")
	}
	data, err = adminGetOrders(adminToken, 0, 10, "created", u1Data.ID, t)
	expectedTotal := 1
	if data.Total != expectedTotal {
		t.Errorf("expected to get %v orders but got %v", expectedTotal, data.Total)
	}

}

func TestVerifyOrder(t *testing.T) {
	officeToken := GetContect("officeToken").(string)

	o1 := GetContect("o1").(models.OrderDto)
	res := Fetch("POST", fmt.Sprintf("order/verify/%v", o1.Id), "", nil, t)
	if x := strings.Contains(string(res), "unathorized"); !x {
		t.Logf("expected to get an error but didn't")
	}

	res = Fetch("POST", fmt.Sprintf("order/verify/%v", o1.Id), officeToken, nil, t)
	response := new(models.ResponseDTO)
	json.Unmarshal(res, response)

	expectedStatus := http.StatusAccepted
	if response.StatusCode != expectedStatus {
		t.Logf("expected to get status %v but got %v", expectedStatus, response.StatusCode)
	}

	u1Token := GetContect("u1Token").(string)
	o1Data, err := getSingleOrder(u1Token, o1.Id, t)
	if err != nil {
		t.Error("expected u1 to be able to get his order")
	}

	expectedOrderStatus := "verified"
	if o1Data.Status != expectedOrderStatus {
		t.Logf("expected status of %v but got %v", expectedOrderStatus, o1Data.Status)
	}

}

func TestPackageOrder(t *testing.T) {
	stockToken := GetContect("stockToken").(string)

	o1 := GetContect("o1").(models.OrderDto)
	res := Fetch("POST", fmt.Sprintf("order/package/%v", o1.Id), "", nil, t)
	if x := strings.Contains(string(res), "unathorized"); !x {
		t.Logf("expected to get an error but didn't")
	}

	res = Fetch("POST", fmt.Sprintf("order/package/%v", o1.Id), stockToken, nil, t)
	response := new(models.ResponseDTO)
	json.Unmarshal(res, response)

	expectedStatus := http.StatusAccepted
	if response.StatusCode != expectedStatus {
		t.Logf("expected to get status %v but got %v", expectedStatus, response.StatusCode)
	}

	u1Token := GetContect("u1Token").(string)
	o1Data, err := getSingleOrder(u1Token, o1.Id, t)
	if err != nil {
		t.Error("expected u1 to be able to get his order")
	}

	expectedOrderStatus := "packaged"
	if o1Data.Status != expectedOrderStatus {
		t.Logf("expected status of %v but got %v", expectedOrderStatus, o1Data.Status)
	}
}

func TestSendOrder(t *testing.T) {
	officeToken := GetContect("officeToken").(string)

	o1 := GetContect("o1").(models.OrderDto)
	res := Fetch("POST", fmt.Sprintf("order/send/%v", o1.Id), "", nil, t)
	if x := strings.Contains(string(res), "unathorized"); !x {
		t.Logf("expected to get an error but didn't")
	}

	res = Fetch("POST", fmt.Sprintf("order/send/%v", o1.Id), officeToken, nil, t)
	response := new(models.ResponseDTO)
	json.Unmarshal(res, response)

	expectedStatus := http.StatusAccepted
	if response.StatusCode != expectedStatus {
		t.Logf("expected to get status %v but got %v", expectedStatus, response.StatusCode)
	}

	u1Token := GetContect("u1Token").(string)
	o1Data, err := getSingleOrder(u1Token, o1.Id, t)
	if err != nil {
		t.Error("expected u1 to be able to get his order")
	}

	expectedOrderStatus := "sent"
	if o1Data.Status != expectedOrderStatus {
		t.Logf("expected status of %v but got %v", expectedOrderStatus, o1Data.Status)
	}
}

func TestDeliverOrder(t *testing.T) {
	shiperToken := GetContect("shiperToken").(string)
	u1Token := GetContect("u1Token").(string)
	o1 := GetContect("o1").(models.OrderDto)

	res := Fetch("GET", fmt.Sprintf("order/shipment-code/%v", o1.Id), u1Token, nil, t)
	t.Log(string(res))

	shipmentCode := ""
	json.Unmarshal(res, &shipmentCode)

	deliverData := models.DeliverOrderDto{
		Code: shipmentCode,
	}
	bodyString, err := json.Marshal(deliverData)
	if err != nil {
		t.Error("this one is on code")
	}
	body := []byte(bodyString)
	res = Fetch("POST", fmt.Sprintf("order/deliver/%v", o1.Id), shiperToken, bytes.NewBuffer(body), t)
	response := new(models.ResponseDTO)
	json.Unmarshal(res, response)
	expectedStatus := http.StatusAccepted
	if response.StatusCode != expectedStatus {
		t.Errorf("expected status of %v but got %v", expectedStatus, response.StatusCode)
	}

	o1Data, err := getSingleOrder(u1Token, o1.Id, t)
	if err != nil {
		t.Error("expected u1 to be able to get his order")
	}

	expectedOrderStatus := "delivered"
	if o1Data.Status != expectedOrderStatus {
		t.Logf("expected status of %v but got %v", expectedOrderStatus, o1Data.Status)
	}

}
