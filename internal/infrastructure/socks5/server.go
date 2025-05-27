package socks5

import (
	"fmt"
	"github.com/things-go/go-socks5"
	"log/slog"
	"net"
	"proxy-server-with-tg-admin/internal/infrastructure/adblock"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
	"sync"
)

type userIdKey string

type Server struct {
	*socks5.Server
	wg     sync.WaitGroup
	stopCh chan struct{}
	logger *slog.Logger
}

func New(statisticTracker *statistic.Tracker, adBlock *adblock.Adblock, authenticator *auth.Authenticator, logger *slog.Logger) *Server {
	srv := socks5.NewServer(
		socks5.WithAuthMethods([]socks5.Authenticator{&UserPassAuthenticator{credentials: &CredentialStore{authenticator: authenticator, logger: logger}}}),
		socks5.WithLogger(&Logger{logger: logger}),
		socks5.WithRule(&Rule{adBlock: adBlock, logger: logger}),
		socks5.WithDialAndRequest(dialAndRequest(statisticTracker, logger)),
		socks5.WithDial(dial(statisticTracker, logger)),
		socks5.WithResolver(DNSResolver{adBlock: adBlock}),
	)

	return &Server{
		Server: srv,
		stopCh: make(chan struct{}),
		logger: logger,
	}
}

func (s *Server) ListenAndServe(network, addr string) error {
	const op = "socks5.ListenAndServe"

	l, err := net.Listen(network, addr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return s.Serve(l)
}

func (s *Server) Serve(l net.Listener) error {
	const op = "socks5.Serve"

	defer l.Close()

	for {
		select {
		case <-s.stopCh:
			return nil
		default:
			conn, err := l.Accept()
			if err != nil {
				return fmt.Errorf("%s: failed to accept connection: %w", op, err)
			}

			errCh := make(chan error)

			go func() {
				errCh <- s.ServeConn(conn)
			}()

			s.wg.Add(1)
			go func() {
				defer s.wg.Done()

				select {
				case <-s.stopCh:
					s.logger.Debug("Socks5 server stopping serve connection")

					if err := conn.Close(); err != nil {
						s.logger.Warn("Socks5 server stopping serve connection", "conn.Close", err)
					}
				case err = <-errCh:
					if err != nil {
						s.logger.Error(op, "serve connection", err)
					}
				}
			}()
		}
	}
}

func (s *Server) Shutdown() {
	s.logger.Debug("Socks5 server shutting down")
	close(s.stopCh)

	waitConnectionsCh := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(waitConnectionsCh)
	}()

	<-waitConnectionsCh
	s.logger.Debug("Socks5 server all connections closed")
}
