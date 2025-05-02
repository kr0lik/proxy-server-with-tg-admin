package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"time"
)

func (s *Storage) AddStat(userId uint32, bytesIn, bytesOut uint64) error {
	const op = "storage.sqlite.AddStat"

	smtp, err := s.db.Prepare(`
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
updated = CURRENT_TIMESTAMP`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := smtp.Exec(userId, bytesIn, bytesOut, bytesIn, bytesOut); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetStatistic(username string) (*entity.UserStat, error) {
	const op = "storage.sqlite.GetStatistic"

	userId, err := s.GetUserId(username)
	if err != nil {
		if errors.Is(err, UserNotFound) {
			return nil, UserNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	smtp, err := s.db.Prepare("SELECT traffic_in_day, traffic_out_day, traffic_in_total, traffic_out_total, days_active, updated FROM user_stat WHERE user_id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userStat := &entity.UserStat{}
	userStat.UserID = userId
	userStat.Updated = time.Time{}

	err = smtp.QueryRow(userId).Scan(&userStat.TrafficInDay, &userStat.TrafficOutDay, &userStat.TrafficInTotal, &userStat.TrafficOutTotal, &userStat.DaysActive, &userStat.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userStat, nil
}

func (s *Storage) DeleteUserStat(username string) error {
	const op = "storage.sqlite.DeleteUserStat"

	smtp, err := s.db.Prepare("DELETE FROM user_stat WHERE username = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = smtp.Exec(username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
