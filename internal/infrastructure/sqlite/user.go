package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"proxy-server-with-tg-admin/internal/entity"
	"strings"
	"time"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")

type scanRow interface {
	Scan(dest ...any) error
}

func (s *Storage) CreateUser(username, password string) (uint32, error) {
	const op = "storage.sqlite.CreateUser"

	res, err := s.db.Exec("INSERT INTO  user(username, password) VALUES(?, ?)", username, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, ErrUserExists
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if id < 0 || id > math.MaxUint32 {
		return 0, fmt.Errorf("%s: id out of range %d", op, id)
	}

	return uint32(id), nil
}

func (s *Storage) GetUser(username string) (*entity.User, error) {
	const op = "storage.sqlite.GetUser"

	row := s.db.QueryRow("SELECT id, username, password, active, ttl, updated FROM user WHERE username = ?", username)
	user, err := s.getEntity(row)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) getUserId(username string) (uint32, error) {
	const op = "storage.sqlite.getUserId"
	var id uint32

	err := s.db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, ErrUserNotFound
		}

		return id, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) ListUsers() ([]*entity.User, error) {
	const op = "storage.sqlite.ListUsers"

	userCount, _ := s.countUsers()

	list := make([]*entity.User, 0, userCount)

	rows, err := s.db.Query("SELECT id, username, password, active, ttl, updated FROM user")
	if err != nil {
		return list, fmt.Errorf("%s: %w", op, err)
	}

	if rows.Err() != nil {
		return list, fmt.Errorf("%s: %w", op, rows.Err())
	}

	for rows.Next() {
		user, err := s.getEntity(rows)
		if err == nil {
			list = append(list, user)
		} else {
			return list, fmt.Errorf("%s: %w", op, err)
		}
	}

	return list, nil
}

func (s *Storage) ActivateUser(username string) error {
	const op = "storage.sqlite.ActivateUser"

	_, err := s.db.Exec("UPDATE user SET active = true, updated = CURRENT_TIMESTAMP WHERE username = ?", username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeactivateUser(username string) error {
	const op = "storage.sqlite.DeactivateUser"

	_, err := s.db.Exec("UPDATE user SET active = false, updated = CURRENT_TIMESTAMP WHERE username = ?", username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdatePassword(username, password string) error {
	const op = "storage.sqlite.UpdatePassword"

	_, err := s.db.Exec("UPDATE user SET password = ?, updated = CURRENT_TIMESTAMP WHERE username = ?", password, username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateTtl(username string, ttl time.Time) error {
	const op = "storage.sqlite.UpdateTtl"

	ttlToUpdate := ttl.Unix()

	if ttlToUpdate < 1 {
		ttlToUpdate = 0
	}

	_, err := s.db.Exec("UPDATE user SET ttl = ?, updated = CURRENT_TIMESTAMP WHERE username = ?", ttlToUpdate, username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUser(username string) error {
	const op = "storage.sqlite.DeleteUser"

	_, err := s.db.Exec("DELETE FROM user WHERE username = ?", username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RenameUser(username, usernameTo string) error {
	const op = "storage.sqlite.RenameUser"

	_, err := s.db.Exec("UPDATE user SET username = ? WHERE username = ?", usernameTo, username)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUserExists
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) getEntity(row scanRow) (*entity.User, error) {
	const op = "storage.sqlite.getEntity"
	var ttl int64

	user := &entity.User{}

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Active, &ttl, &user.Updated)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.Ttl = toTime(ttl)

	return user, nil
}

func (s *Storage) countUsers() (int, error) {
	const op = "storage.sqlite.countUsers"

	count := 0

	err := s.db.QueryRow("SELECT COUNT(id) FROM user").Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return count, ErrUserNotFound
		}

		return count, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
