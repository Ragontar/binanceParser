package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method, http.MethodOptions).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"AssetHistoryGET",
		strings.ToUpper("GET"),
		"/asset/history/{assetName}",
		AssetHistoryGET,
	},

	Route{
		"AssetAddPUT",
		strings.ToUpper("PUT"),
		"/asset/add",
		AssetAddPUT,
	},
}
