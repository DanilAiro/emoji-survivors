package gameservice

import (
	"context"
	"emoji-survivors/game-service/config"
	"emoji-survivors/game-service/db"
	"emoji-survivors/game-service/handlers"
	"emoji-survivors/game-service/initializers"
	"emoji-survivors/game-service/repository"
	"emoji-survivors/game-service/game"
	"emoji-survivors/game-service/ws"
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

	scoreRepo := repository.NewScoreRepository(conn)
	hub := ws.NewHub()
	loop := game.NewGameLoop(hub, scoreRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go loop.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws/connect", handlers.WSConnectHandler(hub, cfg.Secret))

	addr := ":" + cfg.Port
	log.Printf("game-service запущен на %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("сервер завершился с ошибкой: %v", err)
	}
}