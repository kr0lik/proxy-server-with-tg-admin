package socks5

import (
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
	c.logger.Debug("Socks5 dial connection closing", "user", c.UserId, "in", c.BytesRead, "out", c.BytesWritten)

	c.statisticTracker.Track(c.UserId, c.BytesRead, c.BytesWritten)

	return c.Conn.Close()
}

func (c *connection) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	c.BytesRead += uint64(n)
	return n, err
}

func (c *connection) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	c.BytesWritten += uint64(n)
	return n, err
}
