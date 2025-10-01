package socks5

import (
	"context"
	"github.com/things-go/go-socks5"
	"log/slog"
	"proxy-server-with-tg-admin/internal/helper"
	"proxy-server-with-tg-admin/internal/infrastructure/adblock"
)

type Rule struct {
	adBlock *adblock.Adblock
	logger  *slog.Logger
}

func (r Rule) Allow(ctx context.Context, request *socks5.Request) (context.Context, bool) {
	const op = "socks5.Rule.Allow"

	dest := request.DestAddr
	host := dest.FQDN

	if host == "" && dest.IP != nil {
		host = dest.String()
	}

	if r.adBlock.IsMatch(host) {
		r.logger.Warn("Adblock", "ad blocked", host)

		return nil, false
	}

	userIdStr, exist := request.AuthContext.Payload["userId"]
	if !exist {
		r.logger.Error(op, "context", "missing userId")

		return ctx, false
	}

	userId, err := helper.StringToUint32(userIdStr)
	if err != nil {
		r.logger.Error(op, "userId convert", err)

		return ctx, false
	}

	ctx = context.WithValue(ctx, userIdKey("userId"), userId)

	return ctx, true
}
