package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

/*
var drivers []Driver = []Driver{
	Driver{
		Id:        "1",
		Firstname: "David",
		Lastname:  "Noleto",
		Username:  "dbone",
		Password:  "123",
	},
	Driver{
		Id:        "2",
		Firstname: "Fabricio",
		Lastname:  "Martins",
		Username:  "xibungo",
		Password:  "321",
	},
}

var routes []Route = []Route{
	Route{
		Id:        "1",
		Driver:    "2",
		Zipcode:   "94015",
		Numberpkg: "15",
	},
	Route{
		Id:        "2",
		Driver:    "2",
		Zipcode:   "94401",
		Numberpkg: "30",
	},
}
*/
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
	router.HandleFunc("/drivers", DriverRetriveAllEndpoint).Methods("GET")
	router.HandleFunc("/driver/{id}", DriverRetriveEndpoint).Methods("GET")
	router.HandleFunc("/driver/{id}", DriverDeleteEndpoint).Methods("DELETE")
	router.HandleFunc("/driver/{id}", DriverupdateEndpoint).Methods("PUT")
	http.ListenAndServe(":12345", router)

}
