package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
)

const weatherAPIURL = "http://api.weatherapi.com/v1/current.json?key=36a78590564d420288411225250305&q="

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func GetWeatherByCity(city string) (c, f, k float64, err error) {
	tracer := otel.Tracer("weather-service")
	_, span := tracer.Start(context.Background(), "GetWeatherByCity")
	defer span.End()

	resp, err := http.Get(weatherAPIURL + url.QueryEscape(city))
	if err != nil {
		span.RecordError(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("weather API returned %d", resp.StatusCode)
		span.RecordError(err)
		return
	}

	var apiR weatherAPIResponse
	if err = json.NewDecoder(resp.Body).Decode(&apiR); err != nil {
		span.RecordError(err)
		return
	}

	c = apiR.Current.TempC
	f = c*1.8 + 32
	k = c + 273.15
	return
}
