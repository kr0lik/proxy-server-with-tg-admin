package socks5

import (
	"fmt"
	"github.com/things-go/go-socks5"
	"log/slog"
	"net"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
	"sync"
)

type Server struct {
	*socks5.Server
	wg     sync.WaitGroup
	stopCh chan struct{}
	logger *slog.Logger
}

func New(statisticTracker *statistic.Tracker, authenticator *auth.Authenticator, logger *slog.Logger) *Server {
	srv := socks5.NewServer(
		socks5.WithAuthMethods([]socks5.Authenticator{&UserPassAuthenticator{credentials: &CredentialStore{authenticator: authenticator, logger: logger}}}),
		socks5.WithLogger(&Logger{logger: logger}),
		socks5.WithDialAndRequest(dialAndRequest(statisticTracker, logger)),
		socks5.WithDial(dial(logger)),
	)

	return &Server{
		Server: srv,
		stopCh: make(chan struct{}),
		logger: logger,
	}
}

func (s *Server) ListenAndServe(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return fmt.Errorf("socks5.ListenAndServe: %w", err)
	}

	return s.Serve(l)
}

func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		select {
		case <-s.stopCh:
			return nil
		default:
			conn, err := l.Accept()
			if err != nil {
				return fmt.Errorf("socks5.Serve: failed to accept connection: %w", err)
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
						s.logger.Error("Socks5.server", "serve connection", err)
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
