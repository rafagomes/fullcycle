package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const weatherAPIURL = "http://api.weatherapi.com/v1/current.json?key=36a78590564d420288411225250305&q="

// WeatherResponse represents the structure of the weather data response
type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// GetWeatherByCity fetches the current temperature for a given city
func GetWeatherByCity(city string) (float64, float64, float64, error) {
	resp, err := http.Get(weatherAPIURL + url.QueryEscape(city))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, 0, fmt.Errorf("failed to get weather data, status code: %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to decode weather response: %v", err)
	}

	tempC := weatherResponse.Current.TempC
	tempF := CelsiusToFahrenheit(tempC)
	tempK := CelsiusToKelvin(tempC)

	return tempC, tempF, tempK, nil
}

// CelsiusToFahrenheit converts Celsius to Fahrenheit
func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

// CelsiusToKelvin converts Celsius to Kelvin
func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
