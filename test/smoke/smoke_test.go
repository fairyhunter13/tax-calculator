// +build smoke

package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/fairyhunter13/tax-calculator/internal/bill"
	billDelivery "github.com/fairyhunter13/tax-calculator/internal/bill/delivery"

	"github.com/fairyhunter13/tax-calculator/internal/taxobj"

	"github.com/fairyhunter13/tax-calculator/internal/app"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	ini "gopkg.in/ini.v1"
)

var (
	id     int64
	config *app.Config
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	host string
)

const (
	configPath = `/configs/config.ini`
)

func setup(t *testing.T) {
	config = new(app.Config)
	err := ini.MapTo(config, configPath)
	if err != nil {
		t.Fatalf("Error in loading the config: %s", err)
		return
	}
	hostname := os.Getenv("SERVICE_HOST")
	if hostname == "" {
		hostname = "taxcalculator"
	}
	host = "http://" + hostname + config.Server.Port
}

func TestSmoke(t *testing.T) {
	defer cleanup(t)
	setup(t)
	testCreateTaxObject(t)
	testGetBill(t)
}

func testCreateTaxObject(t *testing.T) {
	const requestTest = `
		{
			"name": "KFC Burger",
			"tax_code" : 1,
			"price": 5000
		}
	`
	resp, err := client.Post(host+"/tax", echo.MIMEApplicationJSON, strings.NewReader(requestTest))
	if err != nil {
		t.Fatalf("Error in creating the tax object: %s", err)
		return
	}
	defer resp.Body.Close()
	//Check if the tax object is successfully created.
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	taxObj := new(taxobj.TaxObject)
	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error in reading the response body: %s", err)
		return
	}
	err = json.Unmarshal(byteBody, taxObj)
	if err != nil {
		t.Fatalf("Error in unmarshaling the response body: %s", err)
		return
	}
	if taxObj.ID != 0 {
		id = taxObj.ID
	}
	t.Logf("Tax Object: %+v\n", taxObj)
}

func testGetBill(t *testing.T) {
	expectedResponse := billDelivery.BillResponse{
		Bill: []bill.Bill{
			bill.Bill{
				Name:       "KFC Burger",
				TaxCode:    1,
				Price:      5000,
				Refundable: "Yes",
				Type:       "Food & Beverage",
				Tax:        500,
				Amount:     5500,
			},
		},
		Total: bill.Total{
			PriceSubtotal: 5000,
			TaxSubtotal:   500,
			GrandTotal:    5500,
		},
	}
	resp, err := client.Get(host + "/bill")
	if err != nil {
		t.Fatalf("Error in creating the tax object: %s", err)
		return
	}
	defer resp.Body.Close()
	billResp := new(billDelivery.BillResponse)
	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error in reading the response body: %s", err)
		return
	}
	err = json.Unmarshal(byteBody, billResp)
	if err != nil {
		t.Fatalf("Error in unmarshaling the response body: %s", err)
		return
	}
	t.Logf("Bill List: %+v\n", billResp)
	exist := false
	uniqueCounter := 0
	for _, value := range billResp.Bill {
		if assert.ObjectsAreEqualValues(value, expectedResponse.Bill[0]) && uniqueCounter == 0 {
			uniqueCounter++
			exist = true
			continue
		}
		billResp.Total.PriceSubtotal -= value.Price
		billResp.Total.TaxSubtotal -= value.Tax
		billResp.Total.GrandTotal -= value.Amount
	}
	if !exist {
		t.Fatal("Bill is not exist in the response!")
	}
	assert.Equal(t, expectedResponse.Total.PriceSubtotal, billResp.Total.PriceSubtotal)
	assert.Equal(t, expectedResponse.Total.TaxSubtotal, billResp.Total.TaxSubtotal)
	assert.Equal(t, expectedResponse.Total.GrandTotal, billResp.Total.GrandTotal)
}

func cleanup(t *testing.T) {
	if id != 0 {
		t.Logf("The newly example of the created tax object have id: %d", id)
	}
}
