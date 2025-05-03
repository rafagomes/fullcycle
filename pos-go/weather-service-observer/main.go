package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/trace"

	"weather-service/handlers"
)

func initTracer() func() {
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("failed to initialize zipkin exporter: %v", err)
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(provider)

	return func() {
		_ = provider.Shutdown(context.Background())
	}
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	http.HandleFunc("/weather", handlers.HandleWeatherRequest)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
