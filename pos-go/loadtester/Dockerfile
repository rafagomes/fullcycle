# BUILD
FROM golang:1.23-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod ./

COPY . .
RUN go build -a -ldflags="-s -w" -o loadtester ./cmd/loadtester

# RUN
FROM scratch
COPY --from=builder /app/loadtester /loadtester

ENTRYPOINT ["/loadtester"]