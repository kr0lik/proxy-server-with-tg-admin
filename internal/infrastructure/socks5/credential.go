package socks5

import (
	"fmt"
	"log/slog"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type CredentialStore struct {
	authenticator *auth.Authenticator
	logger        *slog.Logger
}

func (c *CredentialStore) Valid(username, password, userAddr string) bool {
	userId, err := c.authenticator.Authenticate(username, password)
	if err != nil {
		c.logger.Error("Socks5 authentication error", "err", err)

		return false
	}

	return userId != 0
}

func (c *CredentialStore) GetUserId(username, password, userAddr string) (uint32, error) {
	const op = "socks5.CredentialStore.GetUserId"

	userId, err := c.authenticator.Authenticate(username, password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}
