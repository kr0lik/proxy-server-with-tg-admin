package socks5

import (
	"log/slog"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type CredentialStore struct {
	authenticator *auth.Authenticator
	logger        *slog.Logger
}

func (c *CredentialStore) Valid(username, password, userAddr string) bool {
	return c.authenticator.Authenticate(username, password)
}
