package historyManager

import (
	"encoding/json"
	"math"
	"sync"
	"time"
)

type Direction string
const (
	directionUp Direction = "up"
	directionDown Direction = "down"
	directionNothingChanged = "-"
)

var directionsMap = map[string]Direction{
	"up": directionUp,
	"down": directionDown,
	"-": directionNothingChanged,
}

type HistoryManager struct {
	Asset Asset
	EntriesBuffer []HistoryEntry // is used to temp store data before saving to storage

	mu sync.Mutex
	storage HistoryStorage
}

type HistoryEntry struct {
	ID string
	Asset Asset
	price float64
	prevPrice float64
	direction Direction
	perc float64
	date time.Time
}

type Asset struct {
	ID string
	Name string
}

func NewHistoryManager(hs HistoryStorage, a Asset) *HistoryManager {
	hm := HistoryManager{
		Asset: a,
		storage: hs,
	}

	// var err error
	// hm.HistoryEntries, err = hm.storage.Load(a.Name)

	return &hm
}

// Saves buffer to storage
func (hm *HistoryManager) SaveBuffer() error {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	 if err := hm.storage.Save(hm.EntriesBuffer); err != nil {
		return err
	 }

	//  hm.HistoryEntries = append(hm.HistoryEntries, hm.EntriesBuffer...)
	 hm.EntriesBuffer = nil

	return nil
}

func (hm *HistoryManager) GetEntriesAsJSON(limit int, offset int) ([]byte, error) {
	if err := hm.SaveBuffer(); err != nil {
		return nil, err
	}
	hm.mu.Lock()
	defer hm.mu.Unlock()
	entries, err := hm.storage.Load(hm.Asset.ID, limit, offset)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(entries)
	return data, err
}

func (hm *HistoryManager) AddHistoryEntry(he HistoryEntry) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	// Get asset last price
	if len(hm.EntriesBuffer) != 0 {
		he.prevPrice = hm.EntriesBuffer[len(hm.EntriesBuffer)-1].price
	} else {
		prevEntry, err := hm.storage.Load(he.Asset.ID, 1, 0) // get latest history entry
		if err != nil {
			return err
		}
		he.prevPrice = prevEntry[0].price
	}

	if he.prevPrice == 0 {
		he.direction = directionNothingChanged
		he.perc = 0
	}
	if he.prevPrice > he.price {
		he.direction = directionDown
		he.perc = math.Abs((1-(he.price/he.prevPrice))*100)
	}
	if he.prevPrice < he.price {
		he.direction = directionUp
		he.perc = math.Abs((1-(he.price/he.prevPrice))*100)
	}
	he.date = time.Now()

	hm.EntriesBuffer = append(hm.EntriesBuffer, he)


	return nil
}