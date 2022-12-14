package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	prs "github.com/Ragontar/binanceParcer/parser"
)

var Parser *prs.Parser

const timeout = 5 * time.Minute

type AssetAddRequestBody struct {
	Symbol string `json:"symbol,omitempty"`
}

func AssetHistoryGET(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	symbol := r.URL.Query().Get("symbol")

	if limitStr == "" || offsetStr == "" || symbol == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No limit and/or offset and/or symbol provided"))
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) //testing purposes
		return
	}
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) //testing purposes
		return
	}

	hm, ok := Parser.HistoryManagersMap[symbol]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong symbol"))
		return
	}

	responseBody, err := hm.GetEntriesAsJSON(int(limit), int(offset))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // testing
		return
	}

	w.Write(responseBody)
}

func AssetAddPOST(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // testing
		return
	}

	var b AssetAddRequestBody
	json.Unmarshal(body, &b)
	if b.Symbol == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) //testing purposes
		return
	}

	if !ValidateSymbol(b.Symbol) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong symbol"))
		return
	}
	if err := Parser.AddAsset(b.Symbol); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // testing
		return
	}

	w.Write([]byte("added"))
}

func ValidateSymbol(symbol string) bool {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=" + symbol)
	if err != nil {
		return false
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}

	return false
}