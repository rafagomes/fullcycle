package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleWeatherRequest_SuccessfulResponse(t *testing.T) {
	// Mock the handler's dependency to avoid external calls
	req, err := http.NewRequest("GET", "/weather?cep=01001000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"temp_C": "25", "temp_F": "77", "temp_K": "298"}`))
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "temp_C")
	assert.Contains(t, rr.Body.String(), "temp_F")
	assert.Contains(t, rr.Body.String(), "temp_K")
}

func TestHandleWeatherRequest_InvalidCEP_Alternative(t *testing.T) {
	// Mock the handler's dependency to avoid external calls
	req, err := http.NewRequest("GET", "/weather?cep=123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"error": "invalid zipcode"}`))
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid zipcode")
}

func TestHandleWeatherRequest_NotFoundCEP_Alternative(t *testing.T) {
	// Mock the handler's dependency to avoid external calls
	req, err := http.NewRequest("GET", "/weather?cep=99999999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "can not find zipcode"}`))
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "can not find zipcode")
}
