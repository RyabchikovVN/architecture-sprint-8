package main

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("oNwoLQdvJAvRcL89SydqCWCe5ry1jMgq")

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reports", reportsHandler).Methods("GET")
	http.ListenAndServe(":8000", r)
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
		fmt.Fprint(w, "Token ERROR" )
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["roles"] != "prothetic_user" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Need role is nothing")
			return
		}
		fmt.Fprint(w, "OK - WELCOME - " + tokenString)
	} else {		
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Token not valid")
	}	

}