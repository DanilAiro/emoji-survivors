package authservice

import (
	"emoji-survivors/auth-service/config"
	"emoji-survivors/auth-service/db"
	"emoji-survivors/auth-service/initializers"
	"emoji-survivors/auth-service/repository"
	"emoji-survivors/auth-service/handlers"
	"log"
	"net/http"
)

func main() {
	initializers.LoadEnvVariables()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка конфигурации: %v", err)
	}

	conn := db.ConnectDB(cfg.DBUrl)
	defer conn.Close()

	userRepo := repository.NewUserRepository(conn)
	scoreRepo := repository.NewScoreRepository(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", handlers.RegisterHandler(userRepo))
	mux.HandleFunc("POST /api/login", handlers.LoginHandler(userRepo, cfg.Secret))
	mux.HandleFunc("GET /api/stats", handlers.StatsHandler(scoreRepo, cfg.Secret))

	addr := ":" + cfg.Port
	log.Printf("auth-service запущен на %s", addr)
	if err := http.ListenAndServe(addr, handlers.WithCORS(mux)); err != nil {
		log.Fatalf("сервер завершился с ошибкой: %v", err)
	}
}
