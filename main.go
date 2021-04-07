package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var drivers = Driver{
	Id:        "1",
	Firstname: "David",
	Lastname:  "Noleto",
	Username:  "dbone",
	Password:  "123",
}

var routes = Route{
	Id:        "1",
	Driver:    "2",
	Zipcode:   "94015",
	Numberpkg: "15",
}

func RootEndpoint(response http.ResponseWriter, request *http.Request) {
	//defining the response type as json, it could be html, xml ...
	response.Header().Add("content-type", "application/json")
	response.Write([]byte(`{ "message": "Hello API" }`))

}

func main() {
	fmt.Println("Starting the application...")
	//initializing gorilla mux router and serving on port :12345
	router := mux.NewRouter()
	router.HandleFunc("/", RootEndpoint).Methods("GET")
	router.HandleFunc("/register", RegisterEndpoint).Methods("POST")
	router.HandleFunc("/login", LoginEndpoint).Methods("POST")
	router.HandleFunc("/drivers", DriverRetriveAllEndpoint).Methods("GET")
	router.HandleFunc("/driver/{id}", DriverRetriveEndpoint).Methods("GET")
	router.HandleFunc("/driver/{id}", DriverDeleteEndpoint).Methods("DELETE")
	router.HandleFunc("/driver/{id}", DriverupdateEndpoint).Methods("PUT")
	router.HandleFunc("/routes", RouteRetrieveAllEndpoint).Methods("GET")
	router.HandleFunc("/route/{id}", RouteRetrieveEndpoint).Methods("GET")
	router.HandleFunc("/route/{id}", RouteDeleteEndpoint).Methods("DELETE")
	router.HandleFunc("/route/{id}", RouteUpdateEndpoint).Methods("PUT")
	router.HandleFunc("/route", RouteCreateEndpoint).Methods("POST")
	http.ListenAndServe(":12345", router)

}
