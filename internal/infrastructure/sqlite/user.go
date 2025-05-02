package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"strings"
	"time"
)

var UserNotFound = errors.New("user not found")
var UserExists = errors.New("user already exists")

type scanRow interface {
	Scan(dest ...any) error
}

func (s *Storage) CreateUser(username, password string) (uint32, error) {
	const op = "storage.sqlite.CreateUser"

	smtp, err := s.db.Prepare("INSERT INTO  user(username, password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := smtp.Exec(username, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, UserExists
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return uint32(id), nil
}

func (s *Storage) GetUser(username string) (*entity.User, error) {
	const op = "storage.sqlite.GetUser"

	smtp, err := s.db.Prepare("SELECT id, username, password, active, ttl, updated FROM user WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := smtp.QueryRow(username)
	user, err := s.getEntity(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, UserNotFound
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetUserId(username string) (uint32, error) {
	const op = "storage.sqlite.GetUserId"
	var id uint32

	smtp, err := s.db.Prepare("SELECT id FROM user WHERE username = ?")
	if err != nil {
		return id, fmt.Errorf("%s: %w", op, err)
	}

	err = smtp.QueryRow(username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, UserNotFound
		}

		return id, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) ListUsers() []*entity.User {
	const op = "storage.sqlite.ListUsers"
	list := make([]*entity.User, 0, 10)

	rows, err := s.db.Query("SELECT id, username, password, active, ttl, updated FROM user")
	if err != nil {
		return list
	}

	for rows.Next() {
		user, err := s.getEntity(rows)
		if err == nil {
			list = append(list, user)
		}
	}

	return list
}

func (s *Storage) ActivateUser(username string) error {
	const op = "storage.sqlite.ActivateUser"

	smtp, err := s.db.Prepare("UPDATE user SET active = true, updated = CURRENT_TIMESTAMP WHERE username = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = smtp.Exec(username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeactivateUser(username string) error {
	const op = "storage.sqlite.DeactivateUser"

	smtp, err := s.db.Prepare("UPDATE user SET active = false, updated = CURRENT_TIMESTAMP WHERE username = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = smtp.Exec(username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdatePassword(username, password string) error {
	const op = "storage.sqlite.UpdatePassword"

	smtp, err := s.db.Prepare("UPDATE user SET password = ?, updated = CURRENT_TIMESTAMP WHERE username = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = smtp.Exec(password, username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateTtl(username string, ttl time.Time) error {
	const op = "storage.sqlite.UpdateTtl"

	smtp, err := s.db.Prepare("UPDATE user SET ttl = ?, updated = CURRENT_TIMESTAMP WHERE username = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	ttlToUpdate := ttl.Unix()

	if ttlToUpdate < 1 {
		ttlToUpdate = 0
	}

	_, err = smtp.Exec(ttlToUpdate, username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUser(username string) error {
	const op = "storage.sqlite.DeleteUser"

	smtp, err := s.db.Prepare("DELETE FROM user WHERE username = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = smtp.Exec(username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) getEntity(row scanRow) (*entity.User, error) {
	user := &entity.User{}
	var ttl int64

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Active, &ttl, &user.Updated)
	if err != nil {
		return nil, err
	}

	if time.Unix(ttl, 0).Unix() > 86400 {
		user.Ttl = time.Unix(ttl, 0)
	}

	return user, nil
}
