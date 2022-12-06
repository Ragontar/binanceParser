package parser

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
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
}

type AssetStorage interface {
	LoadAssets() ([]historyManager.Asset, error)
	AddAsset(historyManager.Asset) error
}

// Creates new Parser with selected AssetStorage. Fetch Interval can be passed in opts.
// Default: 1 second
func NewParser(as AssetStorage, opts... time.Duration) (*Parser, error) {
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
	var as historyManager.Asset
	if _, ok := p.HistoryManagersMap[symbol]; ok {
		return nil
	}

	assetList, err := p.AssetStorage.LoadAssets()
	if err != nil {
		return err
	}

	for _, asset := range assetList {
		if asset.Name == symbol {
			as = historyManager.Asset{
				ID:   asset.ID,
				Name: symbol,
			}
			p.HistoryManagersMap[symbol] = historyManager.NewHistoryManager(
				historyManager.HistoryStorageDB,
				as,
			)
			return nil
		}
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
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	for _, hm := range p.HistoryManagersMap {
		q.Add("symbols", hm.Asset.Name)
	}
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
	var TickerResponse TickerResponse
	if err := json.Unmarshal(body, &TickerResponse); err != nil {
		return err
	}
	log.Println(TickerResponse)

	for _, tickerEntry := range TickerResponse.Entries {
		hm, ok := p.HistoryManagersMap[tickerEntry.symbol]
		if !ok {
			log.Printf("[WARN]: trying to access missing key %s\n", tickerEntry.symbol)
			continue
		}
		priceFloat, err := strconv.ParseFloat(tickerEntry.price, 64)
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
