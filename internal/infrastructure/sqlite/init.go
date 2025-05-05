package sqlite

import (
	"fmt"
	_ "modernc.org/sqlite"
)

func (s *Storage) init() error {
	stmt, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS user (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username varchar(32) NOT NULL,
    	password varchar(32) NOT NULL,
    	limit_speed INTEGER DEFAULT 0 NOT NULL,
    	limit_traffic INTEGER DEFAULT 0 NOT NULL,
    	limit_type INTEGER DEFAULT 0 NOT NULL,
    	active bool DEFAULT FALSE NOT NULL,
    	ttl TIMESTAMP DEFAULT 0 NOT NULL,
    	updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	); CREATE UNIQUE INDEX IF NOT EXISTS username_uniq_idx ON user(username);`)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("sql.prepare user: %w", err)
	}

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("sql.exec user: %w", err)
	}

	stmt, err = s.db.Prepare(`CREATE TABLE IF NOT EXISTS user_stat (
    	user_id INTEGER PRIMARY KEY,
    	traffic_in_day UNSIGNED BIG INT DEFAULT 0 NOT NULL,
    	traffic_out_day UNSIGNED BIG INT DEFAULT 0 NOT NULL,
    	traffic_in_total UNSIGNED BIG INT DEFAULT 0 NOT NULL,
    	traffic_out_total UNSIGNED BIG INT DEFAULT 0 NOT NULL,
    	days_active INTEGER DEFAULT 0 NOT NULL,
    	updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	)`)
	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("sql.prepare user_stat: %w", err)
	}

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("sql.exec user_stat: %w", err)
	}

	return nil
}
