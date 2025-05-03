# Weather Service

This project is a weather service developed in Go that retrieves weather information based on a Brazilian postal code (CEP). It utilizes external APIs to fetch location data and current weather conditions.

## Features

- Receives a valid 8-digit CEP.
- Retrieves the city name using the viaCEP API.
- Fetches the current temperature in Celsius, Fahrenheit, and Kelvin using the WeatherAPI.
- Returns appropriate HTTP responses based on the request outcome.

## Project Structure

```
pos-go
└── weather-service
    ├── src
    │   ├── main.go                # Entry point of the application
    │   ├── handlers
    │   │   └── weather_handler.go  # Handles weather requests
    │   ├── services
    │   │   ├── cep_service.go      # Interacts with viaCEP API
    │   │   └── weather_service.go   # Interacts with WeatherAPI
    │   ├── utils
    │   │   └── conversions.go       # Temperature conversion functions
    │   └── tests
    │       ├── weather_handler_test.go  # Tests for weather handler
    │       ├── cep_service_test.go      # Tests for CEP service
    │       └── weather_service_test.go   # Tests for weather service
    ├── Dockerfile                     # Docker configuration
    ├── docker-compose.yml              # Docker Compose configuration
    ├── go.mod                          # Go module definition
    ├── go.sum                          # Go module checksums
    └── README.md                       # Project documentation
```

## Setup Instructions

1. Clone the repository:

   ```
   git clone https://github.com/microsoft/vscode-remote-try-go.git
   cd pos-go/weather-service
   ```

2. Build the Docker image:

   ```
   docker build -t weather-service .
   ```

3. Run the application using Docker Compose:
   ```
   docker-compose up
   ```

## Usage

Send a GET request to the endpoint with a valid CEP:

```
GET /weather?cep=<valid_cep>
```

### Example Response

On success:

```
HTTP/1.1 200 OK
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

On invalid CEP format:

```
HTTP/1.1 422 Unprocessable Entity
{
  "message": "invalid zipcode"
}
```

On CEP not found:

```
HTTP/1.1 404 Not Found
{
  "message": "can not find zipcode"
}
```

## Deployment

This application was deployed using Google Cloud Run and can be accessed at the following URL:
[Weather Service on Google Cloud Run](https://weather-service-26048141879.us-central1.run.app/weather?cep=24754210)

## Distributed Tracing with OpenTelemetry and Zipkin

This project implements distributed tracing using OpenTelemetry (OTEL) and Zipkin. The tracing is configured to measure the response times of the following services:

1. **CEP Service**: Measures the time taken to fetch location data based on a CEP (ZIP code).
2. **Weather Service**: Measures the time taken to fetch weather data for a given city.

### Configuration

- **Tracer Initialization**: The tracer is initialized in `main.go` and configured to export spans to a Zipkin server running at `http://localhost:9411`.
- **Instrumentation**:
  - `cep_service.go`: Contains spans to measure the response time of the `GetLocationByCEP` function.
  - `weather_service.go`: Contains spans to measure the response time of the `GetWeatherByCity` function.

### Prerequisites

- Ensure that a Zipkin server is running locally on port 9411.
- Install the required dependencies listed in `go.mod`.

### Running the Application

1. Start the Zipkin server.
2. Run the application using `go run main.go`.
3. Access the Zipkin UI at `http://localhost:9411` to view the traces.
