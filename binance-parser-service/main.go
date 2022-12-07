package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ragontar/binanceParcer/parser"
	"github.com/Ragontar/binanceParcer/server"
)

func main() {
	router := server.NewRouter()

	go func(parser *parser.Parser) {
		for {
			time.Sleep(parser.FetchInterval)
			fmt.Println("--------FETCHING------------")
			err := parser.Fetch()
			if err != nil {
				log.Printf("[FETCH]: %v\n", err)
			}
		}
	}(server.Parser)
	http.ListenAndServe("0.0.0.0:8080", router)
}
