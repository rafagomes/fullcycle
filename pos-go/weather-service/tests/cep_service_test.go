package tests

import (
	"net/http"
	"testing"

	"weather-service/services"
)

func TestGetLocationByCEP(t *testing.T) {
	tests := []struct {
		cep          string
		expectedCity string
		expectedCode int
	}{
		{"01001000", "São Paulo", http.StatusOK},
		{"99999999", "", http.StatusNotFound},
		{"12345678", "", http.StatusUnprocessableEntity},
	}

	for _, test := range tests {
		t.Run(test.cep, func(t *testing.T) {
			city, err := services.GetLocationByCEP(test.cep)
			if test.expectedCode == http.StatusOK {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if city != test.expectedCity {
					t.Errorf("Expected city %s, got %s", test.expectedCity, city)
				}
			} else {
				if err == nil {
					t.Errorf("Expected an error, got none")
				}
			}
		})
	}
}
