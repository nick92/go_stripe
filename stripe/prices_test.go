package stripe_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"trucknav/router"

	"github.com/stretchr/testify/assert"
)

func TestPricesGet(t *testing.T) {
	os.Setenv("STRIPE_KEY", "sk_test_51LOSjkDssEaLCZedvZ8TylNw4aLmu7JWOp4PiLH9usx0fdNBisLQmk4ZmYuPxB4vkUiylbQ0Dgj1u16mdWII4p3b00EvxSRjgQ")
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/prices/get", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
