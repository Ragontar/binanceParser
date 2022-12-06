package main

import (
	"net/http"
	"time"

	"github.com/Ragontar/binanceParcer/historyManager"
	"github.com/Ragontar/binanceParcer/parser"
	"github.com/Ragontar/binanceParcer/server"
)

func main() {
	router := server.NewRouter()
	p, err := parser.NewParser(historyManager.HistoryStorageDB)
	if err != nil {
		panic(err)
	}
	go func (parser *parser.Parser)  {
		for {
			time.Sleep(p.FetchInterval)
			parser.Fetch()
		}
	} (p)
	http.ListenAndServe("0.0.0.0:8080", router)
}
