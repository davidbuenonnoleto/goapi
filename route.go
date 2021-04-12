package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/go-playground/validator.v9"
)

type Route struct {
	Id        string `json:"id,omitempty" validate:"omitempty,uuid"`
	User      string `json:"user,omitempty" validate:"isdefault"`
	Zipcode   string `json:"zipcode,omitempty" validate:"required"`
	Numberpkg string `json:"numberpkg,omitempty" validate:"required"`
	Type      string `json:"type,omitempty"`
}

func RouteRetrieveAllEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var routes []Route
	query := gocb.NewN1qlQuery(`SELECT ` + bucket.Name() + `.* FROM ` + bucket.Name() + ` WHERE type = 'route'`)
	rows, _ := bucket.ExecuteN1qlQuery(query, nil)
	var row Route
	for rows.Next(&row) {
		routes = append(routes, row)
	}
	json.NewEncoder(response).Encode(routes)
}

func RouteRetrieveEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	var route Route
	_, err := bucket.Get(params["id"], &route)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(route)
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
	route.User = token.Id
	route.Type = "article"
	bucket.Insert(route.Id, route, 0)
	json.NewEncoder(response).Encode(route)
}

func RouteDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	token := context.Get(request, "decoded").(CustomJWTClaim)
	query := gocb.NewN1qlQuery(`DELETE FROM ` + bucket.Name() + ` WHERE id = $1 AND user = $2 AND type = 'route'`)
	bucket.ExecuteN1qlQuery(query, []interface{}{params["id"], token.Id})
	response.Write([]byte(`{ "message": "` + params["id"] + `" }`))
}

func RouteUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	var changes Route
	json.NewDecoder(request.Body).Decode(&changes)
	token := context.Get(request, "decoded").(CustomJWTClaim)
	queryStr := `UPDATE ` + bucket.Name() + ` SET type = 'route'`
	if changes.Zipcode != "" {
		queryStr += `, zipcode = $zipcode`
	}
	if changes.Numberpkg != "" {
		queryStr += `, numberpkg = $numberpkg`
	}
	queryStr += ` WHERE id = $id AND user = $user AND type = 'route'`
	query := gocb.NewN1qlQuery(queryStr)
	queryParams := map[string]interface{}{
		"zipcode":   changes.Zipcode,
		"numberpkg": changes.Numberpkg,
		"id":        params["id"],
		"user":      token.Id,
	}
	_, err := bucket.ExecuteN1qlQuery(query, queryParams)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	response.Write([]byte(`{ "message": "` + params["id"] + `" }`))
}
