package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_key")

type Credentials struct {
	Password string `json:"password"`
}

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.StandardClaims
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Декодируем JSON
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil || creds.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получаем пароль из переменной окружения
	storedPassword := os.Getenv("TODO_PASSWORD")
	if storedPassword == "" || creds.Password != storedPassword {
		http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	// Создаем токен
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		PasswordHash: storedPassword,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid || claims.PasswordHash != pass {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
