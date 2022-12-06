package historyManager

import (
	"encoding/json"
	"log"
	"math"
	"sync"
	"time"
)

type Direction string

const (
	directionUp             Direction = "up"
	directionDown           Direction = "down"
	directionNothingChanged           = "-"
)

type HistoryManager struct {
	Asset         Asset
	EntriesBuffer []HistoryEntry // is used to temp store data before saving to storage

	mu                   sync.Mutex
	storage              HistoryStorage
	bufferUnloadInterval time.Duration
}

type HistoryEntry struct {
	ID        string    `json:"id,omitempty"`
	Asset     Asset     `json:"asset,omitempty"`
	Price     float64   `json:"price,omitempty"`
	PrevPrice float64   `json:"prev_price,omitempty"`
	Direction Direction `json:"direction,omitempty"`
	Perc      float64   `json:"perc,omitempty"`
	Date      time.Time `json:"date,omitempty"`
}

type Asset struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Initialize history manager for a specific asset.
// Optional parameter changes bufferUnloadInterval.
// Default: 10 seconds respectively.
func NewHistoryManager(hs HistoryStorage, a Asset, params ...time.Duration) *HistoryManager {
	hm := HistoryManager{
		Asset:   a,
		storage: hs,
	}

	if len(params) < 1 {
		hm.bufferUnloadInterval = 10 * time.Second
	} else {
		hm.bufferUnloadInterval = params[0]
	}

	go hm.StartHistoryBufferProcessor()

	return &hm
}

// Saves buffer to storage
func (hm *HistoryManager) SaveBuffer() error {
	if len(hm.EntriesBuffer) == 0 {
		return nil
	}
	hm.mu.Lock()
	defer hm.mu.Unlock()
	if err := hm.storage.Save(hm.EntriesBuffer); err != nil {
		return err
	}
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
		he.PrevPrice = hm.EntriesBuffer[len(hm.EntriesBuffer)-1].Price
	} else {
		prevEntry, err := hm.storage.Load(he.Asset.ID, 1, 0) // get latest history entry
		if err != nil {
			return err
		}
		if len(prevEntry) == 0 {
			he.PrevPrice = he.Price
		} else {
			he.PrevPrice = prevEntry[0].Price
		}
	}

	if he.PrevPrice == 0 {
		he.Direction = directionNothingChanged
		he.Perc = 0
	}
	if he.PrevPrice > he.Price {
		he.Direction = directionDown
		he.Perc = math.Abs((1 - (he.Price / he.PrevPrice)) * 100)
	}
	if he.PrevPrice < he.Price {
		he.Direction = directionUp
		he.Perc = math.Abs((1 - (he.Price / he.PrevPrice)) * 100)
	}
	he.Date = time.Now()

	hm.EntriesBuffer = append(hm.EntriesBuffer, he)

	return nil
}

func (hm *HistoryManager) StartHistoryBufferProcessor() {
	for {
		time.Sleep(hm.bufferUnloadInterval) // cringe?
		if err := hm.SaveBuffer(); err != nil {
			log.Printf("[Save buffer]: %v\n", err)
		}
	}
}
