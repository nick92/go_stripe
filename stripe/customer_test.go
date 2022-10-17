package stripe_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"trucknav/models"
	"trucknav/router"

	"github.com/stretchr/testify/assert"
)

func TestCustomerGet(t *testing.T) {
	os.Setenv("STRIPE_KEY", "sk_test_51LOSjkDssEaLCZedvZ8TylNw4aLmu7JWOp4PiLH9usx0fdNBisLQmk4ZmYuPxB4vkUiylbQ0Dgj1u16mdWII4p3b00EvxSRjgQ")
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/customer/get?customer_id=cus_MCstCrF3NZMC81", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCustomerAmendDefault(t *testing.T) {
	os.Setenv("STRIPE_KEY", "sk_test_51LOSjkDssEaLCZedvZ8TylNw4aLmu7JWOp4PiLH9usx0fdNBisLQmk4ZmYuPxB4vkUiylbQ0Dgj1u16mdWII4p3b00EvxSRjgQ")
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/customer/set_default?customer_id=cus_MCstCrF3NZMC81&payment_id=card_1LnKAqDssEaLCZedo3wUn7vB", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCustomerAdd(t *testing.T) {
	os.Setenv("STRIPE_KEY", "sk_test_51LOSjkDssEaLCZedvZ8TylNw4aLmu7JWOp4PiLH9usx0fdNBisLQmk4ZmYuPxB4vkUiylbQ0Dgj1u16mdWII4p3b00EvxSRjgQ")
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	pmreq, _ := http.NewRequest(http.MethodPost, "/api/prices/add", nil)
	router.ServeHTTP(w, pmreq)

	var data models.AddPaymentMethodResponse

	json.Unmarshal(w.Body.Bytes(), &data)

	if !assert.Equal(t, http.StatusOK, w.Code) || !assert.NotEmpty(t, data.Complete) {
		assert.Fail(t, "Create payment method failure")
	}

	body := models.AddCustomerRequest{
		Name:            "Bob Jones",
		Email:           "bob@email.com",
		PostCode:        "CH21DF",
		Country:         "GB",
		PaymentMethodId: data.Complete,
		ProdId:          "prod_M6pUxYANJ3WwSH",
	}

	jsonBody, _ := json.Marshal(body)

	ww := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/customer/add", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(ww, req)

	var customerData models.AddCustomerTestResponse

	json.Unmarshal(ww.Body.Bytes(), &customerData)

	assert.Equal(t, http.StatusOK, ww.Code)
	assert.NotEmpty(t, customerData.Complete.CustomerId)
}
