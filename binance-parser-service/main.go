package main

import (
	"net/http"

	"github.com/Ragontar/binanceParcer/server"
)

func main() {
	router := server.NewRouter()
	http.ListenAndServe("0.0.0.0:8080", router)
}
