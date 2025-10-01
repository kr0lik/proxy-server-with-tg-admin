package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
)

const fileName = "server.db"

type Storage struct {
	storePath string
	db        *sql.DB
	logger    *slog.Logger
}

func New(storagePath string, logger *slog.Logger) (*Storage, error) {
	storePath := storagePath + string(os.PathSeparator) + fileName

	s := &Storage{storePath: storePath, logger: logger}
	if err := s.open(); err != nil {
		return nil, fmt.Errorf("sqlite.init: %w", err)
	}

	if err := s.init(); err != nil {
		s.Close()

		return nil, fmt.Errorf("sqlite.init: %w", err)
	}

	return s, nil
}

func (s *Storage) open() error {
	db, err := sql.Open("sqlite3", s.storePath)
	if err != nil {
		return fmt.Errorf("sqlite.open: %w", err)
	}

	_, err = db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA wal_autocheckpoint = 500;
		PRAGMA journal_size_limit = 41943040;
		PRAGMA busy_timeout = 10000;
		PRAGMA locking_mode = NORMAL;
	`)
	if err != nil {
		return fmt.Errorf("sqlite.exec PRAGMA: %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close() {
	if closErr := s.db.Close(); closErr != nil {
		s.logger.Error("Sqlite.bd", "close", closErr)
	}

	s.logger.Debug("Sqlite db closed")
}
