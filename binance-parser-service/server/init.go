package server

import (
	"log"

	"github.com/Ragontar/binanceParcer/historyManager"
	prs "github.com/Ragontar/binanceParcer/parser"
)

const apiUrl = "https://api.binance.com/api/v3/ticker/price"

func init() {
	var err error
	Parser, err = prs.NewParser(historyManager.HistoryStorageDB, apiUrl)
	if err != nil {
		log.Println(err)
		panic("cannot initialize parser")
	}
}