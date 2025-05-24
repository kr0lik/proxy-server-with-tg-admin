package socks5

import (
	"context"
	"errors"
	"fmt"
	"net"
	"proxy-server-with-tg-admin/internal/infrastructure/adblock"
)

var ErrAdBlocked = errors.New("AdBlocked")

type DNSResolver struct {
	adBlock *adblock.Adblock
}

func (d DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	const op = "socks5.DNSResolver.Resolve"

	if d.adBlock.IsMatch(name) {
		return ctx, nil, ErrAdBlocked
	}

	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", op, err)
	}

	return ctx, addr.IP, nil
}
