package parser

type TickerEntry struct {
	Symbol string `json:"symbol,omitempty"`
	Price  string `json:"price,omitempty"`
}

type TickerResponse struct {
	Entries []TickerEntry
}