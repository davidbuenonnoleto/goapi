package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type Driver struct {
	Id        string `json:"id,omitempty" validate:"omitempty,uuid"`
	Firstname string `json:"firstname,omitempty" validate:"required"`
	Lastname  string `json:"lastname,omitempty" validate:"required"`
	Username  string `json:"username,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
}

var drivers []Driver = []Driver{
	{
		Id:        "driver-1",
		Firstname: "Nic",
		Lastname:  "Raboy",
		Username:  "nraboy",
		Password:  "pass",
	},
	{
		Id:        "driver-2",
		Firstname: "Maria",
		Lastname:  "Raboy",
		Username:  "mraboy",
		Password:  "abc123",
	},
}

func RegisterEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var driver Driver

	json.NewDecoder(request.Body).Decode(&driver)
	validate := validator.New()
	err := validate.Struct(driver)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(driver.Password), 10)
	driver.Id = uuid.Must(uuid.NewV4()).String()
	driver.Password = string(hash)
	drivers = append(drivers, driver)
	json.NewEncoder(response).Encode(drivers)
}

func LoginEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var data Driver
	json.NewDecoder(request.Body).Decode(&data)
	validate := validator.New()
	err := validate.StructExcept(data, "Firstname", "Lastname")
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	for _, driver := range drivers {
		if driver.Username == data.Username {
			err := bcrypt.CompareHashAndPassword([]byte(driver.Password), []byte(data.Password))
			if err != nil {
				response.WriteHeader(500)
				response.Write([]byte(`{ "message": "invalid password" }`))
				return
			}
			claims := CustomJWTClaim{
				Id: driver.Id,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Local().Add(time.Hour).Unix(),
					Issuer:    "The Polyglot Developer",
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, _ := token.SignedString(JWT_SECRET)
			response.Write([]byte(`{ "token": "` + tokenString + `" }`))
			return
		}
	}
	response.Write([]byte(`{ "message": "invalid username" }`))
}

func DriverRetrieveAllEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode(drivers)
}

func DriverRetrieveEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	for _, driver := range drivers {
		if driver.Id == params["id"] {
			json.NewEncoder(response).Encode(driver)
			return
		}
	}
	json.NewEncoder(response).Encode(Driver{})
}

func DriverDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	for index, driver := range drivers {
		if driver.Id == params["id"] {
			drivers = append(drivers[:index], drivers[index+1:]...)
			json.NewEncoder(response).Encode(drivers)
			return
		}
	}
	json.NewEncoder(response).Encode(Driver{})
}

func DriverUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	var changes Driver
	json.NewDecoder(request.Body).Decode(&changes)
	validate := validator.New()
	err := validate.StructExcept(changes, "Firstname", "Lastname", "Username", "Password")
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	for index, driver := range drivers {
		if driver.Id == params["id"] {
			if changes.Firstname != "" {
				driver.Firstname = changes.Firstname
			}
			if changes.Lastname != "" {
				driver.Lastname = changes.Lastname
			}
			if changes.Username != "" {
				driver.Username = changes.Username
			}
			if changes.Password != "" {
				err = validate.Var(changes.Password, "gte=4")
				if err != nil {
					response.WriteHeader(500)
					response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
					return
				}
				hash, _ := bcrypt.GenerateFromPassword([]byte(changes.Password), 10)
				driver.Password = string(hash)
			}
			drivers[index] = driver
			json.NewEncoder(response).Encode(drivers)
			return
		}
	}
	json.NewEncoder(response).Encode(Driver{})
}
