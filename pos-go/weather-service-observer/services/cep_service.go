package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
)

const viaCEPURL = "https://viacep.com.br/ws/%s/json/"

var (
	ErrInvalidZipcode = errors.New("invalid zipcode")
	ErrNotFound       = errors.New("cannot find zipcode")
)

var cepRe = regexp.MustCompile(`^\d{5}-?\d{3}$`)

// Location represents the structure of the location data returned by the viaCEP API.
type Location struct {
	City string `json:"localidade"`
}

// GetLocationByCEP validates the CEP, calls ViaCEP, and returns either the city or a sentinel error.
func GetLocationByCEP(cep string) (string, error) {
	tracer := otel.Tracer("cep-service")
	_, span := tracer.Start(context.Background(), "GetLocationByCEP")
	defer span.End()

	if !cepRe.MatchString(cep) {
		return "", ErrInvalidZipcode
	}

	resp, err := http.Get(fmt.Sprintf(viaCEPURL, cep))
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrNotFound
	}

	var loc Location
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		span.RecordError(err)
		return "", err
	}

	if loc.City == "" {
		return "", ErrNotFound
	}
	return loc.City, nil
}
