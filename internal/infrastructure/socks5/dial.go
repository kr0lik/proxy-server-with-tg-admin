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
		userIdStr, exist := request.AuthContext.Payload["userId"]
		if !exist {
			return nil, ErrNoUserId
		}

		userId, err := helper.StringToUint32(userIdStr)
		if err != nil {
			logger.Error(op, "userId convert", err)

			return nil, ErrNoUserId
		}

		conn, err := getConn(userId, network, addr, ctx, statisticTracker, logger)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return conn, nil
	}
}

func dial(statisticTracker *statistic.Tracker, logger *slog.Logger) func(ctx context.Context, network, addr string) (net.Conn, error) {
	const op = "socks5.dial"

	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		userIdAny := ctx.Value(userIdKey("userId"))
		if userIdAny == nil {
			return nil, ErrNoUserId
		}

		userId, ok := userIdAny.(uint32)
		if !ok {
			logger.Error(op, "userId", "is not uint32")

			return nil, ErrNoUserId
		}

		conn, err := getConn(userId, network, addr, ctx, statisticTracker, logger)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return conn, nil
	}
}

func getConn(userId uint32, network, addr string, ctx context.Context, statisticTracker *statistic.Tracker, logger *slog.Logger) (net.Conn, error) {
	const op = "socks5.getConn"

	if userId == 0 {
		return nil, ErrNoUserId
	}

	var d net.Dialer

	conn, err := d.DialContext(ctx, network, addr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &connection{Conn: conn, userId: userId, statisticTracker: statisticTracker, logger: logger}, nil
}
