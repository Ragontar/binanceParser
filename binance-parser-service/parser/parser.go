package parser

import (
	"github.com/Ragontar/binanceParcer/historyManager"
	"github.com/google/uuid"
)

const api = "https://api.binance.com/api/v3/ticker/price"

// ! INITIALIZE THE MAP
type Parser struct {
	HistoryManagersMap map[string]*historyManager.HistoryManager // key = symbol
	AssetStorage       AssetStorage
}

type AssetStorage interface {
	LoadAssets() ([]historyManager.Asset, error)
}

func NewParser(as AssetStorage) (*Parser, error) {
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
	return p, nil
}

func (p *Parser) AddAsset(symbol string) {
	// TODO try to GET price before adding
	// TODO try to get asset_id from DB

	// Asset is not present
	if _, ok := p.HistoryManagersMap[symbol]; !ok {
		as := historyManager.Asset{
			ID:   uuid.NewString(),
			Name: symbol,
		}
		p.HistoryManagersMap[symbol] = historyManager.NewHistoryManager(
			historyManager.HistoryStorageDB,
			as,
		)
	}
}

func (p *Parser) LoadAssetsFromStorage(storage historyManager.HistoryStorage) ([]historyManager.Asset, error) {
	assets := []historyManager.Asset{}

	return
}
