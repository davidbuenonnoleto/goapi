package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

type Route struct {
	Id        string `json:"id,omitempty"`
	Driver    string `json:"driver,omitempty"`
	Zipcode   string `json:"zipcode,omitempty"`
	Numberpkg string `json:"numberpkg,omitempty"`
}

var routes []Route = []Route{
	{
		Id:        "route-1",
		Driver:    "driver-1",
		Zipcode:   "94015",
		Numberpkg: "15",
	},
}

func RouteRetrieveAllEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode(routes)
}

func RouteRetrieveEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	for _, route := range routes {
		if route.Id == params["id"] {
			json.NewEncoder(response).Encode(route)
			return
		}
	}
	json.NewEncoder(response).Encode(Route{})
}

func RouteCreateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var route Route
	json.NewDecoder(request.Body).Decode(&route)
	token := context.Get(request, "decoded").(CustomJWTClaim)
	validate := validator.New()
	err := validate.Struct(route)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	route.Id = uuid.Must(uuid.NewV4()).String()
	route.Driver = token.Id
	routes = append(routes, route)
	json.NewEncoder(response).Encode(routes)
}

func RouteDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	token := context.Get(request, "decoded").(CustomJWTClaim)
	for index, route := range routes {
		if route.Id == params["id"] && route.Driver == token.Id {
			routes = append(routes[:index], routes[index+1:]...)
			json.NewEncoder(response).Encode(routes)
			return
		}
	}
	json.NewEncoder(response).Encode(Route{})
}

func RouteUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	var changes Route
	json.NewDecoder(request.Body).Decode(&changes)
	token := context.Get(request, "decoded").(CustomJWTClaim)
	for index, route := range routes {
		if route.Id == params["id"] && route.Driver == token.Id {
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
	}
	json.NewEncoder(response).Encode(Route{})
}
