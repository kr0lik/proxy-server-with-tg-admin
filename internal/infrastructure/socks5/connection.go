package socks5

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
)

type connection struct {
	net.Conn
	BytesRead        uint64
	BytesWritten     uint64
	UserId           uint32
	statisticTracker *statistic.Tracker
	logger           *slog.Logger
}

func (c *connection) Close() error {
	const op = "socks5.connection.Close"

	if c.BytesRead > 0 || c.BytesWritten > 0 {
		c.logger.Debug("Socks5 dial connection closing", "user", c.UserId, "in", c.BytesRead, "out", c.BytesWritten)

		c.statisticTracker.Track(c.UserId, c.BytesRead, c.BytesWritten)
		c.BytesRead, c.BytesWritten = 0, 0
	}

	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *connection) Read(b []byte) (int, error) {
	const op = "socks5.connection.Read"

	n, err := c.Conn.Read(b)
	if err != nil {
		if err.Error() == io.EOF.Error() {
			return n, nil
		}

		return n, fmt.Errorf("%s: %w (%db)", op, err, n)
	}

	if n < 0 {
		return n, fmt.Errorf("%s: negative byte count: %d", op, n)
	}

	if n > 0 {
		c.BytesRead += uint64(n)
	}

	return n, nil
}

func (c *connection) Write(b []byte) (int, error) {
	const op = "socks5.connection.Write"

	n, err := c.Conn.Write(b)
	if err != nil {
		return n, fmt.Errorf("%s: %w (%db)", op, err, n)
	}

	if n < 0 {
		return n, fmt.Errorf("%s: negative byte count: %d", op, n)
	}

	if n > 0 {
		c.BytesWritten += uint64(n)
	}

	return n, nil
}
