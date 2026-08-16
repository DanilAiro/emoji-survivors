package handlers

import (
	"emoji-survivors/auth-service/repository"
	"encoding/json"
	"net/http"

	"emoji-survivors/shared/jwt"

	"golang.org/x/crypto/bcrypt"
)

type RequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(userRepo *repository.UserRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rd RequestData

		err := json.NewDecoder(r.Body).Decode(&rd)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(rd.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Упс, что-то сломалось", http.StatusInternalServerError)
			return
		}

		_, err = userRepo.Create(r.Context(), rd.Username, string(passwordHash))
		if err != nil {
			http.Error(w, "Упс, что-то сломалось", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"result": "Успешная регистрация!",
		})
	}
}

func LoginHandler(userRepo *repository.UserRepository, JWTSecret string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		defer r.Body.Close()

		var rd RequestData

		err := json.NewDecoder(r.Body).Decode(&rd)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		user, err := userRepo.FindByUsername(r.Context(), rd.Username)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(rd.Password))
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		}

		jwtToken, err := jwt.CreateToken(JWTSecret, rd.Username)
		if err != nil {
			http.Error(w, "Упс, что-то сломалось", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"result": "Успешный логин!",
			"jwt":    jwtToken,
		})
	}
}

func StatsHandler(scoreRepo *repository.ScoreRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
