package server

import (
	"log"

	"github.com/Ragontar/binanceParcer/historyManager"
	prs "github.com/Ragontar/binanceParcer/parser"
)

func init() {
	var err error
	Parser, err = prs.NewParser(historyManager.HistoryStorageDB)
	if err != nil {
		log.Println(err)
		panic("cannot initialize parser")
	}
}