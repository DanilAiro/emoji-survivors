package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	appjwt "emoji-survivors/shared/jwt"

	"emoji-survivors/game-service/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WSConnectHandler(hub *ws.Hub, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "отсутствует токен", http.StatusUnauthorized)
			return
		}

		claims, err := appjwt.VerifyToken(token, jwtSecret)
		if err != nil {
			http.Error(w, "невалидный или истёкший токен", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ошибка апгрейда WebSocket-соединения: %v", err)
			return
		}

		player := ws.NewPlayer(claims.UserID, claims.Username, conn)
		hub.Add(player)
		log.Printf("игрок %s подключился", claims.Username)

		go player.WritePump()
		player.ReadPump(func() {
			hub.Remove(claims.UserID)
			log.Printf("игрок %s отключился", claims.Username)
		})
	}
}
