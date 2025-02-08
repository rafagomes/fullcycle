package main

import (
	"log"
	"net/http"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Não foi encontrado .env, utilizando variáveis de ambiente")
	}

	cfg := config.LoadConfig()

	store, err := store.NewRedisStore(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Erro ao conectar com o Redis: %v", err)
	}

	rateLimiter := limiter.NewRateLimiter(store, cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	handler := rateLimiter.Middleware(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("Servidor rodando na porta 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Erro no servidor: %v", err)
	}
}
