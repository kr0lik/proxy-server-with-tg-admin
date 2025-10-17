package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/usecase/commands"
	"strings"
	"time"
)

const SELECT_USER_FIELDS = "id, username, password, COALESCE(telegram_id, 0) as telegram_id, active, ttl, updated"

var ErrUserNotFound = errors.New("user not found")

type scanRow interface {
	Scan(dest ...any) error
}

func (s *Storage) CreateUser(username, password string) (uint32, error) {
	const op = "storage.sqlite.CreateUser"

	res, err := s.db.Exec("INSERT INTO  user(username, password) VALUES($1, $2)", username, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, commands.ErrUserExists
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

func (s *Storage) GetUserByUsername(username string) (*entity.User, error) {
	const op = "storage.sqlite.GetUserByUsername"

	row := s.db.QueryRow("SELECT "+SELECT_USER_FIELDS+" FROM user WHERE username = $1", username)

	user, err := s.getEntity(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetUserByTelegramId(telegramId int64) (*entity.User, error) {
	const op = "storage.sqlite.GetUserByTelegramId"

	row := s.db.QueryRow("SELECT "+SELECT_USER_FIELDS+" FROM user WHERE telegram_id = $1", telegramId)

	user, err := s.getEntity(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetUserIdByUsername(username string) (uint32, error) {
	const op = "storage.sqlite.getUserIdByUsername"
	var id uint32

	err := s.db.QueryRow("SELECT id FROM user WHERE username =$1", username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, ErrUserNotFound
		}

		return id, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserByInviteToken(token string) (*entity.User, error) {
	const op = "storage.sqlite.GetUserByInviteToken"

	row := s.db.QueryRow("SELECT "+SELECT_USER_FIELDS+" FROM user WHERE invite_token = $1", token)

	user, err := s.getEntity(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) ListUsers() ([]*entity.User, error) {
	const op = "storage.sqlite.ListUsers"

	userCount, _ := s.countUsers()

	list := make([]*entity.User, 0, userCount)

	rows, err := s.db.Query("SELECT " + SELECT_USER_FIELDS + " FROM user")
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

	_, err := s.db.Exec("UPDATE user SET active = true, updated = CURRENT_TIMESTAMP WHERE username = $1", username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeactivateUser(username string) error {
	const op = "storage.sqlite.DeactivateUser"

	_, err := s.db.Exec("UPDATE user SET active = false, updated = CURRENT_TIMESTAMP WHERE username = $1", username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdatePassword(username, password string) error {
	const op = "storage.sqlite.UpdatePassword"

	_, err := s.db.Exec("UPDATE user SET password = $2, updated = CURRENT_TIMESTAMP WHERE username = $1", username, password)
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

	_, err := s.db.Exec("UPDATE user SET ttl = $2, updated = CURRENT_TIMESTAMP WHERE username = $1", username, ttlToUpdate)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateInviteToken(username string, token string) error {
	const op = "storage.sqlite.UpdateInviteToken"

	_, err := s.db.Exec("UPDATE user SET invite_token = $2, updated = CURRENT_TIMESTAMP WHERE username = $1", username, token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RenameUser(username, usernameTo string) error {
	const op = "storage.sqlite.RenameUser"

	_, err := s.db.Exec("UPDATE user SET username = $2, updated = CURRENT_TIMESTAMP WHERE username = $1", username, usernameTo)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return commands.ErrUserExists
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AssignTelegramIdByInviteToken(inviteToken string, telegramId int64) error {
	const op = "storage.sqlite.AssignTelegramIdByInviteToken"

	res, err := s.db.Exec("UPDATE user SET telegram_id = $1, invite_token = null, updated = CURRENT_TIMESTAMP WHERE invite_token = $2", telegramId, inviteToken)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return commands.ErrTelegramIdExists
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	rowsUpdated, err := res.RowsAffected()
	if rowsUpdated == 0 && err == nil {
		return fmt.Errorf("%s: %w", op, commands.ErrNoInviteToken)
	}

	return nil
}

func (s *Storage) DeleteUser(userId uint32) error {
	const op = "storage.sqlite.DeleteUser"

	_, err := s.db.Exec("DELETE FROM user WHERE id = $1", userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) getEntity(row scanRow) (*entity.User, error) {
	const op = "storage.sqlite.getEntity"

	user := &entity.User{}

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.TelegramId, &user.Active, &user.Ttl, &user.Updated)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

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
