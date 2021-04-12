package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id        string `json:"id,omitempty" validate:"omitempty,uuid"`
	Firstname string `json:"firstname,omitempty" validate:"required"`
	Lastname  string `json:"lastname,omitempty" validate:"required"`
	Username  string `json:"username,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required,gte=4"`
	Type      string `json:"type,omitempty"`
}

func RegisterEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Id = uuid.Must(uuid.NewV4()).String()
	user.Password = string(hash)
	user.Type = "user"
	_, err = bucket.Insert(user.Id, user, 0)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func LoginEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var data User
	json.NewDecoder(request.Body).Decode(&data)
	validate := validator.New()
	err := validate.StructExcept(data, "Firstname", "Lastname")
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	query := gocb.NewN1qlQuery(`SELECT ` + bucket.Name() + `.* FROM ` + bucket.Name() + ` WHERE username = $1`)
	rows, _ := bucket.ExecuteN1qlQuery(query, []interface{}{data.Username})
	var row User
	err = rows.One(&row)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if row.Id == "" {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "incorrect username" }`))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(data.Password))
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "invalid password" }`))
		return
	}
	claims := CustomJWTClaim{
		Id: row.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour).Unix(),
			Issuer:    "Beyon Monkeys",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JWT_SECRET)
	response.Write([]byte(`{ "token": "` + tokenString + `" }`))
}

func UserRetrieveAllEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var users []User
	query := gocb.NewN1qlQuery(`SELECT ` + bucket.Name() + `.* FROM ` + bucket.Name() + ` WHERE type = 'manager'`)
	rows, _ := bucket.ExecuteN1qlQuery(query, nil)
	var row User
	for rows.Next(&row) {
		users = append(users, row)
	}
	json.NewEncoder(response).Encode(users)
}

func UserRetrieveEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	var user User
	_, err := bucket.Get(params["id"], &user)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func UserDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	_, err := bucket.Remove(params["id"], 0)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	response.Write([]byte(`{ "message": "` + params["id"] + `" }`))
}

func UserUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	var changes User
	json.NewDecoder(request.Body).Decode(&changes)
	validate := validator.New()
	err := validate.StructExcept(changes, "Firstname", "Lastname", "Username", "Password")
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	mutation := bucket.MutateIn(params["id"], 0, 0)
	if changes.Firstname != "" {
		mutation.Upsert("firstname", changes.Firstname, true)
	}
	if changes.Lastname != "" {
		mutation.Upsert("lastname", changes.Lastname, true)
	}
	if changes.Username != "" {
		mutation.Upsert("username", changes.Username, true)
	}
	if changes.Password != "" {
		err = validate.Var(changes.Password, "gte=4")
		if err != nil {
			response.WriteHeader(500)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(changes.Password), 10)
		mutation.Upsert("password", string(hash), true)
	}
	mutation.Execute()
	response.Write([]byte(`{ "message": "` + params["id"] + `" }`))
}
