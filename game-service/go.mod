module emoji-survivors/game-service

go 1.23.3

require (
	github.com/gorilla/websocket v1.5.3
	github.com/joho/godotenv v1.5.1
)

require emoji-survivors/shared/jwt v0.0.0

require github.com/golang-jwt/jwt/v5 v5.3.1 // indirect

replace emoji-survivors/shared/jwt => ../shared/jwt
