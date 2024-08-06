FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN go build -o desafio1 main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/desafio1 /app/desafio1

CMD ["/app/desafio1"]