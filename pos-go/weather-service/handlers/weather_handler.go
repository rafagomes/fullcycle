package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"weather-service/services"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func HandleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	if !isValidCEP(cep) {
		log.Printf("Invalid CEP format: %s", cep)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := services.GetLocationByCEP(cep)
	if err != nil {
		log.Printf("Failed to find location for CEP %s: %v", cep, err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	tempC, tempF, tempK, err := services.GetWeatherByCity(city)
	if err != nil {
		log.Printf("Failed to fetch weather for city %s: %v", city, err)
		http.Error(w, "could not fetch weather", http.StatusInternalServerError)
		return
	}

	response := WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	log.Printf("Successfully fetched weather for city %s: %+v", city, response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func isValidCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	return re.MatchString(cep)
}
