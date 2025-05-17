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
	const op = "socks5.dialAndRequest"

	return func(ctx context.Context, network, addr string, request *socks5.Request) (net.Conn, error) {
		var d net.Dialer

		userIdStr, exist := request.AuthContext.Payload["userId"]
		if !exist {
			return nil, ErrNoUserId
		}

		userId, err := helper.StringToUint32(userIdStr)
		if err != nil {
			logger.Error(op, "userId convert", err)

			return nil, ErrNoUserId
		}

		conn, err := d.DialContext(ctx, network, addr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return &connection{Conn: conn, userId: userId, statisticTracker: statisticTracker, logger: logger}, nil
	}
}

func dial(statisticTracker *statistic.Tracker, logger *slog.Logger) func(ctx context.Context, network, addr string) (net.Conn, error) {
	const op = "socks5.dial"

	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		userId, ok := ctx.Value(userIdKey("userId")).(uint32)
		if !ok {
			logger.Error(op, "userId", "is not uint32")

			return nil, ErrNoUserId
		}

		if userId == 0 {
			logger.Error(op, "context", "missing userId")

			return nil, ErrNoUserId
		}

		var d net.Dialer

		conn, err := d.DialContext(ctx, network, addr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return &connection{Conn: conn, userId: userId, statisticTracker: statisticTracker, logger: logger}, nil
	}
}
