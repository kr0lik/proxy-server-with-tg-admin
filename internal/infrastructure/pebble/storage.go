package pebble

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	jsoniter "github.com/json-iterator/go"
	"log/slog"
	_ "modernc.org/sqlite"
)

var json = jsoniter.ConfigFastest

type Storage struct {
	db     *pebble.DB
	logger *slog.Logger
}

func New(storagePath string, logger *slog.Logger) (*Storage, error) {
	db, err := pebble.Open(storagePath, &pebble.Options{})
	if err != nil {
		return nil, fmt.Errorf("pebble.open: %w", err)
	}

	return &Storage{db: db, logger: logger}, nil
}

func (s *Storage) Close() {
	if closErr := s.db.Close(); closErr != nil {
		s.logger.Error("pebble.close", "err", closErr)
	}

	s.logger.Debug("Pebble db closed")
}
