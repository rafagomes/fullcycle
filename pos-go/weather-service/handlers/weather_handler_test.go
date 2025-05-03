package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleWeatherRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather", func(c *gin.Context) {
		HandleWeatherRequest(c.Writer, c.Request)
	})

	req, _ := http.NewRequest("GET", "/weather?cep=24754210", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "temp_C")
	assert.Contains(t, w.Body.String(), "temp_F")
	assert.Contains(t, w.Body.String(), "temp_K")
}

func TestHandleWeatherRequest_InvalidCEP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather", func(c *gin.Context) {
		HandleWeatherRequest(c.Writer, c.Request)
	})

	req, _ := http.NewRequest("GET", "/weather?cep=1234567", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "invalid zipcode", w.Body.String())
}

func TestHandleWeatherRequest_NotFoundCEP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather", func(c *gin.Context) {
		HandleWeatherRequest(c.Writer, c.Request)
	})

	req, _ := http.NewRequest("GET", "/weather?cep=99999999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "cannot find zipcode", w.Body.String())
}
