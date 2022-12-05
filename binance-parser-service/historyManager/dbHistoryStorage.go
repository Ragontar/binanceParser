package historyManager

import (
	"context"
	"fmt"
	"time"

	"github.com/Ragontar/binanceParcer/env"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HistoryStorage interface {
	Save([]HistoryEntry) error
	Load(assetName string, limit int, offset int) ([]HistoryEntry, error)
}

type DBHistoryStorage struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewDBHistoryStorage() (*DBHistoryStorage, error) {
	s := &DBHistoryStorage{}
	s.timeout = 5 * time.Minute
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
