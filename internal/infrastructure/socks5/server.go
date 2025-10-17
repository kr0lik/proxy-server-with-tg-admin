package socks5

import (
	"errors"
	"fmt"
	"github.com/things-go/go-socks5"
	"io"
	"log/slog"
	"proxy-server-with-tg-admin/internal/infrastructure/adblock"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
	"strings"
)

type userIdKey string

type Server struct {
	*socks5.Server
}

func New(statisticTracker *statistic.Tracker, adBlock *adblock.Adblock, authenticator *auth.Authenticator, logger *slog.Logger) *Server {
	srv := socks5.NewServer(
		socks5.WithAuthMethods([]socks5.Authenticator{&UserPassAuthenticator{credentials: &CredentialStore{authenticator: authenticator, logger: logger}}}),
		socks5.WithRule(&Rule{adBlock: adBlock, logger: logger}),
		socks5.WithDialAndRequest(dialAndRequest(statisticTracker, logger)),
		socks5.WithDial(dial(statisticTracker, logger)),
		socks5.WithResolver(DNSResolver{adBlock: adBlock}),
		socks5.WithLogger(&logAdapter{logger: logger}),
	)

	return &Server{
		Server: srv,
	}
}

type logAdapter struct {
	logger *slog.Logger
}

func (l *logAdapter) Errorf(format string, args ...interface{}) {
	if len(args) > 0 {
		errStr := fmt.Sprintf(format, args...)

		for _, arg := range args {
			if err, ok := arg.(error); ok && isExpectedError(err) {
				l.logger.Debug(errStr)

				return
			}
		}

		if isExpectedErrorString(errStr) {
			l.logger.Debug(errStr)
		} else {
			l.logger.Error(errStr)
		}
	}
}

func isExpectedError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, io.EOF)
}

func isExpectedErrorString(msg string) bool {
	if msg == "" {
		return false
	}

	return strings.Contains(msg, "EOF") ||
		strings.Contains(msg, "client want to used addr")
}
