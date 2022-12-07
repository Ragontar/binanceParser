package historyManager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ragontar/binanceParcer/env"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HistoryStorage interface {
	Save([]HistoryEntry) error
	Load(assetID string, limit int, offset int) ([]HistoryEntry, error)
}

type DBHistoryStorage struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewDBHistoryStorage() (*DBHistoryStorage, error) {
	s := &DBHistoryStorage{}
	s.timeout = 30 * time.Second
	err := s.init()
	return s, err
}

func (s *DBHistoryStorage) init() error {
	if s.db != nil {
		return nil
	}
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		env.DB_USER,
		env.DB_PASSWORD,
		env.DB_ADDR,
		env.DB_DATABASE,
	)

	var err error
	s.db, err = pgxpool.Connect(context.TODO(), dsn)

	return err
}

func (dbs *DBHistoryStorage) Save(entries []HistoryEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout)
	defer cancel()
	var queryString = INSERT_HISTORY_ENTRIES
	for _, he := range entries {
		queryString += fmt.Sprintf(
			"('%s', '%s', %v, '%s', %v, '%s'),",
			he.ID,
			he.Asset.ID,
			he.Price,
			string(he.Direction),
			he.Perc,
			he.Date.Format(time.RFC3339),
		)
	}
	queryString = queryString[0 : len(queryString)-1]
	fmt.Println(queryString)

	rows, err := dbs.db.Query(ctx, queryString)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}

func (dbs *DBHistoryStorage) Load(assetID string, limit int, offset int) ([]HistoryEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout)
	defer cancel()

	rows, err := dbs.db.Query(ctx, SELECT_HISTORY_ENTRIES_BY_ASSET_ID, assetID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]HistoryEntry, 0, limit)

	for rows.Next() {
		var e HistoryEntry
		err := rows.Scan(&e.ID, &e.Asset.ID, &e.Price, &e.Direction, &e.Perc, &e.Date)
		if err != nil {
			log.Printf("[SCAN]: %v\n", err)
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (dbs *DBHistoryStorage) LoadAssets() ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbs.timeout)
	defer cancel()

	rows, err := dbs.db.Query(ctx, SELECT_ASSETS)
	if err != nil {
		return nil, err
	}

	assets := []Asset{}
	for rows.Next() {
		var a Asset
		err := rows.Scan(&a.ID, &a.Name)
		if err != nil {
			log.Printf("[SCAN]: %v\n", err)
		}
		defer rows.Close()
		assets = append(assets, a)
	}

	return assets, nil
}

func (dbs *DBHistoryStorage) AddAsset(a Asset) error {
	ctx, canel := context.WithTimeout(context.Background(), dbs.timeout)
	defer canel()
	rows, err := dbs.db.Query(ctx, INSERT_ASSET, a.ID, a.Name)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}
