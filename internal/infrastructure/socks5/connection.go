package socks5

import (
	"fmt"
	"log/slog"
	"net"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
	"sync/atomic"
)

type connection struct {
	net.Conn
	bytesRead        uint64
	bytesWritten     uint64
	userId           uint32
	statisticTracker *statistic.Tracker
	logger           *slog.Logger
}

func (c *connection) Close() error {
	const op = "socks5.connection.Close"

	bytesRead := atomic.SwapUint64(&c.bytesRead, 0)
	bytesWritten := atomic.SwapUint64(&c.bytesWritten, 0)

	if bytesRead > 0 || bytesWritten > 0 {
		c.logger.Debug("Socks5 connection closing", "user", c.userId, "in", bytesRead, "out", bytesWritten)

		// in: bytesRead - from target to user
		// out: bytesWritten - from user to target
		c.statisticTracker.Track(c.userId, bytesRead, bytesWritten)
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
		return n, fmt.Errorf("%s: %w (%db)", op, err, n)
	}

	if n < 0 {
		return n, fmt.Errorf("%s: negative byte count: %d", op, n)
	}

	if n > 0 {
		atomic.AddUint64(&c.bytesRead, uint64(n))
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
		atomic.AddUint64(&c.bytesWritten, uint64(n))
	}

	return n, nil
}
