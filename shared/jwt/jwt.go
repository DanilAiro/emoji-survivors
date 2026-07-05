package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(email string) (string, error) {
	// Create claims with standard and custom fields
	claims := jwt.MapClaims{
		"sub": email,                                 // Subject
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expiration
		"iat": time.Now().Unix(),                     // Issued At
	}

	// Create token with HS256 algorithm and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(secretKey)
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what you expect (HMAC in this case)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}