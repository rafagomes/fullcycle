package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const weatherAPIURL = "http://api.weatherapi.com/v1/current.json?key=36a78590564d420288411225250305&q="

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func GetWeatherByCity(city string) (c, f, k float64, err error) {
	resp, err := http.Get(weatherAPIURL + url.QueryEscape(city))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("weather API returned %d", resp.StatusCode)
		return
	}

	var apiR weatherAPIResponse
	if err = json.NewDecoder(resp.Body).Decode(&apiR); err != nil {
		return
	}

	c = apiR.Current.TempC
	f = c*1.8 + 32
	k = c + 273.15
	return
}
