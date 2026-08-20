package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTTL = time.Hour * 24
)

var (
	ErrInvalidToker = errors.New("невалидный токен")
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CreateToken(userID int, username, secretKey string) (string, error) {
	claims := Claims{
		UserID: userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with HS256 algorithm and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(secretKey, tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what you expect (HMAC in this case)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToker
	}

	return claims, nil
}
