package main

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// CORSMiddleware добавляет необходимые заголовки для поддержки CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                            // Разрешаем все источники (для разработки)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")          // Разрешаем методы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Разрешаем заголовки

		if r.Method == http.MethodOptions {
			// Обработка preflight-запросов
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r) // Передаем управление следующему обработчику
	})
}

var mySigningKey = []byte("oNwoLQdvJAvRcL89SydqCWCe5ry1jMgq")

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reports", reportsHandler).Methods("GET")
	// Оборачиваем маршрутизатор в middleware CORS
	http.ListenAndServe(":8000", CORSMiddleware(r))
}

func reportsHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Token ERROR!")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["roles"] != "prothetic_user" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Target role not found!")
			return
		}
		fmt.Fprint(w, "OK - WELCOME!")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Token not valid!")
	}

}
