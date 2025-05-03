FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app

RUN go build -o weather-service .

FROM gcr.io/distroless/base

COPY --from=builder /app/weather-service /weather-service

CMD ["./weather-service"]