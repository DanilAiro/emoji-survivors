module emoji-survivors/auth-service

go 1.27

require (
	emoji-survivors/shared/jwt v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.12.3
	golang.org/x/crypto v0.48.0
)

require github.com/golang-jwt/jwt/v5 v5.3.1 // indirect

replace emoji-survivors/shared/jwt => ../shared/jwt
