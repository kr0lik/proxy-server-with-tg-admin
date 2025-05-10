package socks5

import (
	"context"
	"errors"
	"fmt"
	"github.com/things-go/go-socks5"
	"log/slog"
	"net"
	"proxy-server-with-tg-admin/internal/helper"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
)

var ErrNoUserId = errors.New("no user id")

func dialAndRequest(statisticTracker *statistic.Tracker, logger *slog.Logger) func(ctx context.Context, network, addr string, request *socks5.Request) (net.Conn, error) {
	return func(ctx context.Context, network, addr string, request *socks5.Request) (net.Conn, error) {
		var d net.Dialer

		userIdStr, ok := request.AuthContext.Payload["userId"]
		if !ok {
			return nil, ErrNoUserId
		}

		userId, err := helper.StringToUint32(userIdStr)
		if err != nil {
			logger.Error("Socks5.dial", "get user id err", err)

			return nil, ErrNoUserId
		}

		conn, err := d.DialContext(ctx, network, addr)
		if err != nil {
			return nil, fmt.Errorf("dial: %w", err)
		}

		return &connection{Conn: conn, userId: userId, statisticTracker: statisticTracker, logger: logger}, nil
	}
}

func dial(logger *slog.Logger) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var d net.Dialer

		logger.Warn("Simple dial without traffic tracking")

		return d.DialContext(ctx, network, addr)
	}
}
