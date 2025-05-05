package pebble

import (
	"errors"
	"fmt"
	"github.com/cockroachdb/pebble"
	"log/slog"
	"proxy-server-with-tg-admin/internal/entity"
)

var UserNotFound = errors.New("user not found")
var UserExists = errors.New("user already exists")

type UserRepository struct {
	db     *pebble.DB
	logger *slog.Logger
}

func (s *UserRepository) CreateUser(username, password string) error {
	const op = "storage.pebble.CreateUser"

	user := entity.User{Username: username, Password: password}

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.db.Set([]byte(user.Username), data, pebble.Sync); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *UserRepository) GetUser(username string) (*entity.User, error) {
	const op = "storage.pebble.GetUser"

	val, closer, err := s.db.Get([]byte(username))
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	user := &entity.User{}

	err = json.Unmarshal(val, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, err
}

func (s *UserRepository) ListUsers() ([]*entity.User, error) {
	const op = "storage.pebble.ListUsers"

	iter, err := s.db.NewIter(nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer iter.Close()

	users := make([]*entity.User, 0, 10)

	for iter.First(); iter.Valid(); iter.Next() {
		val, err := iter.ValueAndErr()
		if err != nil {
			continue
		}

		user := &entity.User{}

		if err := json.Unmarshal(val, &user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserRepository) SaveUser(user *entity.User) error {
	const op = "storage.pebble.SaveUser"

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.db.Set([]byte(user.Username), data, pebble.Sync); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
