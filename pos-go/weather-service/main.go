package main

import (
	"log"
	"net/http"

	"weather-service/handlers"
)

func main() {
	http.HandleFunc("/weather", handlers.HandleWeatherRequest)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
