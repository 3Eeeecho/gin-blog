package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwaggerRoute(t *testing.T) {
	router := InitRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Swagger")
}
