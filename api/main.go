package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Realm struct {
	PublicKey string `json:"public_key"`
}

var publicKey []byte
var realmPK string

func initRSAPublicKey() {
	resp, err := http.Get("http://keycloak:8080/realms/reports-realm")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var realm Realm
	err = json.Unmarshal(body, &realm)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	realmPK = "\n-----BEGIN PUBLIC KEY-----\n" + realm.PublicKey + "\n-----END PUBLIC KEY-----\n"
	publicKey = []byte("\n-----BEGIN PUBLIC KEY-----\n" + realm.PublicKey + "\n-----END PUBLIC KEY-----\n")
}

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

//var publicKey = []byte("\n-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzEd0L5pqlOVLuyrfWOVDmqn4MvsX8jn+k2Ea/Tsi7xxYysGQaCAvGa3TMdzTw69knyg9VUJTxPSxp2kJ4BRm5jVgEQb6v+44P2gcaEUjPpBgDyGYa2wHuZntjEbl63rF0lJjpjFO2HMw6wh4cOlQjHHnLAKBqw7TD3bXDDUWY2A16viB54vdeLXT+QqvaRBm4pZFRIwIg0V+SVo8yIuvcVdKDphq13d4P0Hf/EdeB4Reg2FKmHdRZwzCXDRqND/0EhC9OftEi53d/K2MPVK4JwH7Z8YrRyZtV+jraEXpaePzO6jDpDyNoLxaN0uAv6al9QY/6UVdmHHlaKtV4A12tQIDAQAB\n-----END PUBLIC KEY-----\n")

func reportsHandler(w http.ResponseWriter, r *http.Request) {
	initRSAPublicKey()

	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	var pubKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "ERROR: "+err.Error()+"- pk: "+realmPK)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	/*

	 */
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized - err: "+err.Error()+"- pk: "+realmPK)
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
