package parser

type TickerEntry struct {
	symbol string `json:"symbol,omitempty"`
	price  string `json:"price,omitempty"`
}

type TickerResponse struct {
	Entries []TickerEntry
}