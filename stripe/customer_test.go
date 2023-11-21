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
	"trucknav/stripe"

	"github.com/stretchr/testify/assert"
)

func TestCustomerGet(t *testing.T) {
	os.Setenv("STRIPE_KEY", "")
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/customer/get?customer_id=cus_MCstCrF3NZMC81", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCustomerAmendDefault(t *testing.T) {
	os.Setenv("STRIPE_KEY", "")
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/customer/set_default?customer_id=cus_MCstCrF3NZMC81&payment_id=card_1LnKAqDssEaLCZedo3wUn7vB", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCustomerAdd(t *testing.T) {
	key := ""
	os.Setenv("STRIPE_KEY", key)
	router := router.SetupRouter()

	card := stripe.CreatePaymentMethod(key)

	body := models.AddCustomerRequest{
		Name:            "Bob Jones",
		Email:           "bob@email.com",
		PostCode:        "CH21DF",
		Country:         "GB",
		PaymentMethodId: card,
		ProdId:          "prod_M6pUxYANJ3WwSH",
	}

	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/customer/add", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	var customerData models.AddCustomerTestResponse

	json.Unmarshal(w.Body.Bytes(), &customerData)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, customerData.Complete.CustomerId)
}

func TestCustomerPaymentAttachDetach(t *testing.T) {
	key := ""
	os.Setenv("STRIPE_KEY", key)
	router := router.SetupRouter()

	card := stripe.CreatePaymentMethod(key)

	body := models.AddCustomerPaymentMethodRequest{
		PaymentMethodId: card,
		CustomerId:      "cus_MCstCrF3NZMC81",
	}

	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/customer/add_card", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodDelete, "/api/customer/delete_card?payment_id="+card, bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, http.StatusOK, w2.Code)
}
