package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"weather-service/services"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// writeError writes a plain-text body (no trailing newline) with the given status code.
func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(msg))
}

func HandleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	city, err := services.GetLocationByCEP(r.URL.Query().Get("cep"))
	if err != nil {
		log.Printf("CEP error: %v", err)
		switch err {
		case services.ErrInvalidZipcode:
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		case services.ErrNotFound:
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "unexpected error")
		}
		return
	}

	tempC, tempF, tempK, err := services.GetWeatherByCity(city)
	if err != nil {
		log.Printf("Weather API error for %q: %v", city, err)
		writeError(w, http.StatusInternalServerError, "could not fetch weather")
		return
	}

	resp := WeatherResponse{TempC: tempC, TempF: tempF, TempK: tempK}
	log.Printf("Weather for %q: %+v", city, resp)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
