package main

import (
	"emoji-survivors/auth-service/config"
	"emoji-survivors/auth-service/db"
	"emoji-survivors/auth-service/handlers"
	"emoji-survivors/auth-service/initializers"
	"emoji-survivors/auth-service/repository"
	"log"
	"net/http"
)

func main() {
	initializers.LoadEnvVariables()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка конфигурации: %v", err)
	}

	conn, err := db.ConnectDB(cfg.DBUrl)
	if err == db.ErrNoConnect || err == db.ErrNoPing {
		log.Fatal(err.Error())
	}
	if err != nil {
		log.Fatalf("ошибка подключения к базе: %v", err)
	}
	defer conn.Close()

	userRepo := repository.NewUserRepository(conn)
	scoreRepo := repository.NewScoreRepository(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", handlers.RegisterHandler(userRepo))
	mux.HandleFunc("POST /api/login", handlers.LoginHandler(userRepo, cfg.Secret))
	mux.HandleFunc("GET /api/stats", handlers.StatsHandler(scoreRepo))

	addr := ":" + cfg.Port
	log.Printf("auth-service запущен на %s", addr)
	if err := http.ListenAndServe(addr, handlers.WithCORS(mux)); err != nil {
		log.Fatalf("сервер завершился с ошибкой: %v", err)
	}
}
