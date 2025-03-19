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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reports", reportsHandler).Methods("GET")
	http.ListenAndServe(":8000", CORSMiddleware(r))
}

var publicKey = []byte("\n-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1mxs5cw/bOGQ+PHLFeIwDAUzGhKLDYbNJaqoHjvV/GUCAp/1FF/TDIiEN2fuLLeAgpPZKZSEG+bxXJ7vXXxSNvUjY1etmRvHbhlMN+rJYTQtWLZxsvM/MfYi0b+l20oeGtCBwUa2CKWFssnxBM5L3Ex7bmTKzetivOrFP25ztCd6Bu5RXaZ/KJgJlFoHqw8GkGp6iMb1bXcyOSlndJWYprHfa5XjbBsNmXwZ6y8AM3V8Qb+dzDeA60qmX5RlWQOp80W50QKh7BOBud0JuLewZE4JKiWUsM5Csisc/XZbYQSpWZTgGduqi1taYDOZHSPA/yJUtIqdXiyCwo5P0oIqKwIDAQAB\n-----END PUBLIC KEY-----\n")

func reportsHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "ERROR: "+err.Error())
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized - err: "+err.Error())
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["realm_access"].(map[string]interface{})["roles"].([]interface{})[0] != "prothetic_user" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		fmt.Fprint(w, "Welcome!")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized")
	}
}
