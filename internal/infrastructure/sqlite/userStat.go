package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/usecase/commands"
	"time"
)

func (s *Storage) AddStat(userId uint32, bytesIn, bytesOut uint64) error {
	const op = "storage.sqlite.AddStat"

	_, err := s.db.Exec(`
INSERT INTO  user_stat(user_id, traffic_in_day, traffic_out_day, traffic_in_total, traffic_out_total, days_active, updated) VALUES(?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(user_id) DO UPDATE SET
traffic_in_total=traffic_in_total + excluded.traffic_in_total,
traffic_out_total=traffic_out_total + excluded.traffic_out_total,
traffic_in_day=CASE
	WHEN strftime('%d', updated) <> strftime('%d', excluded.updated)
	THEN 0
	ELSE traffic_in_day + excluded.traffic_in_day
END,
traffic_out_day=CASE
	WHEN strftime('%d', updated) <> strftime('%d', excluded.updated)
	THEN 0
	ELSE traffic_out_day + excluded.traffic_out_day
END,
days_active=CASE
	WHEN strftime('%d', updated) <> strftime('%d', excluded.updated)
	THEN days_active + 1
	ELSE days_active
END,
updated = CURRENT_TIMESTAMP`, userId, bytesIn, bytesOut, bytesIn, bytesOut)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetStatistic(username string) (*entity.UserStat, error) {
	const op = "storage.sqlite.GetStatistic"

	userId, err := s.getUserId(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userStat := &entity.UserStat{}
	userStat.UserID = userId
	userStat.Updated = time.Time{}

	err = s.db.QueryRow("SELECT traffic_in_day, traffic_out_day, traffic_in_total, traffic_out_total, days_active, updated FROM user_stat WHERE user_id = ?", userId).
		Scan(&userStat.TrafficInDay, &userStat.TrafficOutDay, &userStat.TrafficInTotal, &userStat.TrafficOutTotal, &userStat.DaysActive, &userStat.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userStat, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userStat, nil
}

func (s *Storage) DeleteUserStat(username string) error {
	const op = "storage.sqlite.DeleteUserStat"

	userId, err := s.getUserId(username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("DELETE FROM user_stat WHERE user_id = ?", userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ListUsersWithStat() ([]*commands.UsersWithStatDto, error) {
	const op = "storage.sqlite.ListUsersWithStat"

	userCount, _ := s.countUsers()

	list := make([]*commands.UsersWithStatDto, 0, userCount)

	rows, err := s.db.Query("SELECT u.username, u.active, u.ttl, COALESCE(us.traffic_in_total, 0), COALESCE(us.traffic_out_total, 0), COALESCE(us.days_active, 0), COALESCE(strftime('%s', us.updated), 0) FROM user u LEFT JOIN user_stat us ON u.id = us.user_id")
	defer rows.Close()

	if err != nil {
		return list, fmt.Errorf("%s: %w", op, err)
	}

	if rows.Err() != nil {
		return list, fmt.Errorf("%s: %w", op, rows.Err())
	}

	for rows.Next() {
		var ttl, updated int64
		dto := &commands.UsersWithStatDto{}

		err := rows.Scan(&dto.Username, &dto.Active, &ttl, &dto.TotalIn, &dto.TotalOut, &dto.DyesActive, &updated)
		if err == nil {
			dto.Ttl = toTime(ttl)
			dto.LastActive = toTime(updated)

			list = append(list, dto)
		} else {
			return list, fmt.Errorf("%s: %w", op, err)
		}
	}

	return list, nil
}

func (s *Storage) DeleteUserWithStat(username string) error {
	const op = "storage.sqlite.DeleteUserWithStat"

	if err := s.DeleteUserStat(username); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.DeleteUser(username); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
