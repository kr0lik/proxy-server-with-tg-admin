package socks5

import (
	"context"
	"fmt"
	"github.com/things-go/go-socks5"
	"log/slog"
	"net"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
)

type CustomLogger struct {
	logger *slog.Logger
}

func (cl CustomLogger) Errorf(format string, v ...interface{}) {
	cl.logger.Error("socks5.server", "internal", fmt.Sprintf(format, v...))
}

// TODO
func GetServer(statisticTracker *statistic.Tracker, authenticator *auth.Authenticator, logger *slog.Logger) *socks5.Server {
	return socks5.NewServer(
		socks5.WithCredential(&CredentialStore{authenticator: authenticator, logger: logger}),
		socks5.WithLogger(&CustomLogger{logger: logger}),
		socks5.WithDialAndRequest(func(ctx context.Context, network, addr string, request *socks5.Request) (net.Conn, error) {
			var d net.Dialer

			conn, err := d.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}

			var userId uint32

			// TODO
			username, ok := request.AuthContext.Payload["username"]
			if ok {
				password, ok := request.AuthContext.Payload["password"]
				if ok {
					userId, err = authenticator.GetUserId(username, password)
					if err != nil {
						return nil, err
					}

					return &connection{Conn: conn, UserId: userId, statisticTracker: statisticTracker, logger: logger}, nil
				}
			}

			return nil, fmt.Errorf("bad credentials")
		}),
	)
}
