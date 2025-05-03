package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"weather-service/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleWeatherRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather/:cep", func(c *gin.Context) {
		handlers.HandleWeatherRequest(c.Writer, c.Request)
	})

	req, _ := http.NewRequest("GET", "/weather/01001000", nil)
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
	router.GET("/weather/:cep", func(c *gin.Context) {
		handlers.HandleWeatherRequest(c.Writer, c.Request)
	})

	req, _ := http.NewRequest("GET", "/weather/1234567", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "invalid zipcode", w.Body.String())
}

func TestHandleWeatherRequest_NotFoundCEP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather/:cep", func(c *gin.Context) {
		handlers.HandleWeatherRequest(c.Writer, c.Request)
	})

	req, _ := http.NewRequest("GET", "/weather/99999999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "can not find zipcode", w.Body.String())
}
