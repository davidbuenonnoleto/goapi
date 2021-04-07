package main

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type Route struct {
	Id        string `json:"id,omitempty"`
	Driver    string `json:"driver,omitempty"`
	Zipcode   string `json:"zipcode,omitempty"`
	Numberpkg string `json:"numberpkg,omitempty"`
}

func RouteRetrieveAllEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode(routes)
}

func RouteRetrieveEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	/*params := mux.Vars(request)
	for _, route := range routes {
		if route.Id == params["id"] {
			json.NewEncoder(response).Encode(route)
			return
		}
	}*/
	json.NewEncoder(response).Encode(Route{})
}

func RouteCreateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var route Route
	json.NewDecoder(request.Body).Decode(&route)
	route.Id = uuid.Must(uuid.NewV4()).String()
	//routes = append(routes, route)
	json.NewEncoder(response).Encode(routes)
}

func RouteDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	/*params := mux.Vars(request)
	for index, route := range routes {
		if route.Id == params["id"] {
			routes = append(routes[:index], routes[index+1:]...)
			json.NewEncoder(response).Encode(routes)
			return
		}
	}*/
	json.NewEncoder(response).Encode(Route{})
}

func RouteUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	/*params := mux.Vars(request)
	var changes Route
	json.NewDecoder(request.Body).Decode(&changes)
	for index, route := range routes {
		if route.Id == params["id"] {
			if changes.Zipcode != "" {
				route.Zipcode = changes.Zipcode
			}
			if changes.Numberpkg != "" {
				route.Numberpkg = changes.Numberpkg
			}
			routes[index] = route
			json.NewEncoder(response).Encode(routes)
			return
		}
	}*/
	json.NewEncoder(response).Encode(Route{})
}
