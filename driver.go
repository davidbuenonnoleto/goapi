package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Driver struct {
	Id        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}

func DriverRetriveAllEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode("teste driver")
}

func DriverRetriveEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	for _, driver := range drivers {
		if driver.Id == params["id"] {
			json.NewEncoder(response).Encode(drivers)
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
			//learn more about this line of code
			drivers = append(drivers[:index], drivers[index+1:]...)
			json.NewEncoder(response).Encode(drivers)
			return
		}
	}
	json.NewEncoder(response).Encode(Driver{})
}

func DriverupdateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	var changes Driver
	json.NewDecoder(request.Body).Decode(&changes)
	for index, author := range drivers {
		if author.Id == params["id"] {
			if changes.Firstname != "" {
				author.Firstname = changes.Firstname
			}
			if changes.Lastname != "" {
				author.Lastname = changes.Lastname
			}
			if changes.Username != "" {
				author.Username = changes.Username
			}
			if changes.Password != "" {
				author.Password = changes.Password
			}
			drivers[index] = author
			json.NewEncoder(response).Encode(drivers)
			return
		}
	}
	json.NewEncoder(response).Encode(Driver{})

}
