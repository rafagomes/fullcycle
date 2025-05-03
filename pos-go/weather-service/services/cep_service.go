package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const viaCEPURL = "https://viacep.com.br/ws/%s/json/"

// Location represents the structure of the location data returned by the viaCEP API.
type Location struct {
	City string `json:"localidade"`
}

// GetLocationByCEP retrieves the city name based on the provided CEP.
func GetLocationByCEP(cep string) (string, error) {
	if len(cep) != 8 {
		return "", errors.New("unprocessable entity: invalid zipcode")
	}

	resp, err := http.Get(fmt.Sprintf(viaCEPURL, cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("can not find zipcode")
	}

	var location Location
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return "", err
	}

	return location.City, nil
}
