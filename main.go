package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

type CustomJWTClaim struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

var JWT_SECRET []byte = []byte("jsonwebtokensecretkey")

func ValidateJWT(t string) (interface{}, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return JWT_SECRET, nil
	})
	if err != nil {
		return nil, errors.New(`{ "message": "` + err.Error() + `" }`)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var tokenData CustomJWTClaim
		mapstructure.Decode(claims, &tokenData)
		return tokenData, nil
	} else {
		return nil, errors.New(`{ "message": "invalid token" }`)
	}
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		driverizationHeader := request.Header.Get("driverization")
		if driverizationHeader != "" {
			bearerToken := strings.Split(driverizationHeader, " ")
			if len(bearerToken) == 2 {
				decoded, err := ValidateJWT(bearerToken[1])
				if err != nil {
					response.Header().Add("content-type", "application/json")
					response.WriteHeader(500)
					response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
					return
				}
				context.Set(request, "decoded", decoded)
				next(response, request)
			}
		} else {
			response.Header().Add("content-type", "application/json")
			response.WriteHeader(500)
			response.Write([]byte(`{ "message": "auth header is required" }`))
			return
		}
	})
}

func RootEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Write([]byte(`{ "message": "Hello World" }`))
}

func main() {
	fmt.Println("Starting the application...")
	r := mux.NewRouter()
	r.HandleFunc("/", RootEndpoint).Methods("GET")
	r.HandleFunc("/register", RegisterEndpoint).Methods("POST")
	r.HandleFunc("/login", LoginEndpoint).Methods("POST")
	r.HandleFunc("/drivers", DriverRetrieveAllEndpoint).Methods("GET")
	r.HandleFunc("/driver/{id}", DriverRetrieveEndpoint).Methods("GET")
	r.HandleFunc("/driver/{id}", DriverDeleteEndpoint).Methods("DELETE")
	r.HandleFunc("/driver/{id}", DriverUpdateEndpoint).Methods("PUT")
	r.HandleFunc("/routes", RouteRetrieveAllEndpoint).Methods("GET")
	r.HandleFunc("/route/{id}", RouteRetrieveEndpoint).Methods("GET")
	r.HandleFunc("/route/{id}", ValidateMiddleware(RouteDeleteEndpoint)).Methods("DELETE")
	r.HandleFunc("/route/{id}", ValidateMiddleware(RouteUpdateEndpoint)).Methods("PUT")
	r.HandleFunc("/route", ValidateMiddleware(RouteCreateEndpoint)).Methods("POST")
	http.ListenAndServe(":12345", r)
}
