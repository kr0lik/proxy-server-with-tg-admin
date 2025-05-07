package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	_ "modernc.org/sqlite"
	"os"
)

type Storage struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(storagePath, filename string, logger *slog.Logger) (*Storage, error) {
	db, err := sql.Open("sqlite", storagePath+string(os.PathSeparator)+filename+"?_busy_timeout=1000")
	if err != nil {
		return nil, fmt.Errorf("sqlite.open: %w", err)
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, fmt.Errorf("sqlite.exec PRAGMA journal_mode=WAL: %w", err)
	}

	s := &Storage{db: db, logger: logger}

	if err = s.init(); err != nil {
		s.Close()

		return nil, fmt.Errorf("sqlite.init: %w", err)
	}

	return s, nil
}

func (s *Storage) Close() {
	if closErr := s.db.Close(); closErr != nil {
		s.logger.Error("Sqlite.bd", "close", closErr)
	}

	s.logger.Debug("Sqlite db closed")
}
