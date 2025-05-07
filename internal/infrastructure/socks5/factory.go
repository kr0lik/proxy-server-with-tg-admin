package socks5

import (
	"github.com/things-go/go-socks5"
	"log/slog"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
)

// Deprecated
func MakeServer(statisticTracker *statistic.Tracker, authenticator *auth.Authenticator, logger *slog.Logger) *socks5.Server {
	return socks5.NewServer(
		socks5.WithAuthMethods([]socks5.Authenticator{&UserPassAuthenticator{credentials: &CredentialStore{authenticator: authenticator, logger: logger}}}),
		socks5.WithLogger(&Logger{logger: logger}),
		socks5.WithDialAndRequest(dialAndRequest(statisticTracker, logger)),
	)
}
