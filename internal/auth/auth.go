package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_password")

// Структура для хранения данных о токене
type Claims struct {
	Hash string `json:"hash"`
	jwt.StandardClaims
}

// Обработчик для аутентификации
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if request.Password != pass {
		http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	// Создание JWT-токена
	hash := "some_hash_based_on_password" // Замените на реальный хэш пароля
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		Hash: hash,
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

	// Установка куки с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	// Возврат токена
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Middleware для проверки аутентификации
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtToken string
		cookie, err := r.Cookie("token")
		if err == nil {
			jwtToken = cookie.Value
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid || claims.Hash != "some_hash_based_on_password" { // Замените на реальный хэш пароля
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}
