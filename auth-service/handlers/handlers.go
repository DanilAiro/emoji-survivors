package handlers

import (
	"emoji-survivors/auth-service/models"
	"emoji-survivors/auth-service/repository"
	"encoding/json"
	"net/http"

	"emoji-survivors/shared/jwt"

	"golang.org/x/crypto/bcrypt"
)

type RequestData struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResultResponse struct {
	Result string `json:"result"`
}

type LoginResponse struct {
	ResultResponse
	Token string `json:"token"`
}

type StatisticsResponse struct {
	ResultResponse
	Statistics []models.ScoreWithUsername `json:"statistics"`
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
		json.NewEncoder(w).Encode(ResultResponse{
			Result: "Успешная регистрация!",
		})
	}
}

func LoginHandler(userRepo *repository.UserRepository, secret string) func(http.ResponseWriter, *http.Request) {
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

		token, err := jwt.CreateToken(rd.UserID, rd.Username, secret)
		if err != nil {
			http.Error(w, "Упс, что-то сломалось", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginResponse{
			ResultResponse: ResultResponse{
				Result: "Успешный логин!",
			},
			Token: token,
		})
	}
}

func StatsHandler(scoreRepo *repository.ScoreRepository, secret string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "отсутствует токен", http.StatusUnauthorized)
			return
		}

		_, err := jwt.VerifyToken(token, secret)
		if err != nil {
			http.Error(w, "невалидный или истёкший токен", http.StatusUnauthorized)
			return
		}

		statistics, err := scoreRepo.GetTop10(r.Context())
		if err != nil {
			http.Error(w, "Упс, что-то сломалось", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StatisticsResponse{
			ResultResponse: ResultResponse{
				Result: "Успешный запрос статистики!",
			},
			Statistics: statistics,
		})
	}
}
