package pebble

import (
	"errors"
	"fmt"
	"github.com/cockroachdb/pebble"
	"log/slog"
	"proxy-server-with-tg-admin/internal/entity"
)

type UserStatRepository struct {
	db     *pebble.DB
	logger *slog.Logger
}

func (s *Storage) AddUserStat(username string, bytesIn, bytesOut uint64) error {
	const op = "storage.pebble.AddUserStat"

	return nil
}

func (s *Storage) GetUserStat(username string) (*entity.UserStat, error) {
	const op = "storage.pebble.GetUserStat"

	userStat := &entity.UserStat{}

	val, closer, err := s.db.Get([]byte(username + "-stat"))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return userStat, nil
		}

		return userStat, err
	}
	defer closer.Close()

	err = json.Unmarshal(val, &userStat)
	if err != nil {
		return userStat, fmt.Errorf("%s: %w", op, err)
	}

	return userStat, err
}
