package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Ragontar/binanceParcer/historyManager"
	"github.com/google/uuid"
)

const api = "https://api.binance.com/api/v3/ticker/price"

// ! INITIALIZE THE MAP
type Parser struct {
	HistoryManagersMap map[string]*historyManager.HistoryManager // key = symbol
	AssetStorage       AssetStorage

	FetchInterval time.Duration

	mu sync.Mutex
}

type AssetStorage interface {
	LoadAssets() ([]historyManager.Asset, error)
	AddAsset(historyManager.Asset) error
}

// Creates new Parser with selected AssetStorage. Fetch Interval can be passed in opts.
// Default: 1 second
func NewParser(as AssetStorage, opts ...time.Duration) (*Parser, error) {
	p := &Parser{
		AssetStorage:       as,
		HistoryManagersMap: make(map[string]*historyManager.HistoryManager),
	}
	assets, err := p.AssetStorage.LoadAssets()
	if err != nil {
		return nil, err
	}
	for _, a := range assets {
		p.HistoryManagersMap[a.Name] = historyManager.NewHistoryManager(historyManager.HistoryStorageDB, a)
	}

	if len(opts) == 0 {
		p.FetchInterval = time.Second
	} else {
		p.FetchInterval = opts[0]
	}

	return p, nil
}

func (p *Parser) AddAsset(symbol string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var as historyManager.Asset
	if _, ok := p.HistoryManagersMap[symbol]; ok {
		return nil
	}

	// Asset is not present in DB
	as = historyManager.Asset{
		ID:   uuid.NewString(),
		Name: symbol,
	}
	if err := p.AssetStorage.AddAsset(as); err != nil {
		return err
	}
	p.HistoryManagersMap[symbol] = historyManager.NewHistoryManager(
		historyManager.HistoryStorageDB,
		as,
	)

	return nil
}

// Gets data from Binance API, appends new entries to history buffers of active HistoryManagers
func (p *Parser) Fetch() error {
	if len(p.HistoryManagersMap) == 0 {
		log.Println("HistoryManagersMap is empty...")
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	queryValue := "["
	for _, hm := range p.HistoryManagersMap {
		queryValue += fmt.Sprintf("\"%s\",", hm.Asset.Name)
	}
	queryValue = queryValue[:len(queryValue)-1]
	queryValue += "]"
	q.Add("symbols", queryValue)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// var TickerResponse TickerResponse
	TickerResponse := []TickerEntry{}
	if err := json.Unmarshal(body, &TickerResponse); err != nil {
		return err
	}
	log.Println(TickerResponse)

	for _, tickerEntry := range TickerResponse {
		hm, ok := p.HistoryManagersMap[tickerEntry.Symbol]
		if !ok {
			log.Printf("[WARN]: trying to access missing key %s\n", tickerEntry.Symbol)
			continue
		}
		priceFloat, err := strconv.ParseFloat(tickerEntry.Price, 64)
		if err != nil {
			log.Printf("parsing float: %s\n", err)
			continue
		}
		he := historyManager.HistoryEntry{
			ID:    uuid.NewString(),
			Asset: hm.Asset,
			Price: priceFloat,
		}

		if err := hm.AddHistoryEntry(he); err != nil {
			log.Printf("adding history entry: %s\n", err)
		}
	}

	return nil
}
