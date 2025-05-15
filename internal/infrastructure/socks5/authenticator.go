package socks5

import (
	"errors"
	"fmt"
	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
	"io"
	"proxy-server-with-tg-admin/internal/helper"
	"proxy-server-with-tg-admin/internal/infrastructure/sqlite"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type UserPassAuthenticator struct {
	credentials *CredentialStore
}

func (a UserPassAuthenticator) GetCode() uint8 { return statute.MethodUserPassAuth }

func (a UserPassAuthenticator) Authenticate(reader io.Reader, writer io.Writer, userAddr string) (*socks5.AuthContext, error) {
	const op = "socks5.Authenticate"

	// reply the client to use user/pass auth
	if _, err := writer.Write([]byte{statute.VersionSocks5, statute.MethodUserPassAuth}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// get user and user's password
	nup, err := statute.ParseUserPassRequest(reader)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	username := string(nup.User)
	password := string(nup.Pass)

	// Verify
	userId, err := a.credentials.GetUserId(username, password, userAddr)
	if err != nil {
		if _, err := writer.Write([]byte{statute.UserPassAuthVersion, statute.AuthFailure}); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if errors.Is(err, auth.ErrUserPassword) || errors.Is(err, sqlite.ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w (%s:%s %s)", op, err, username, password, userAddr)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if userId == 0 {
		if _, err := writer.Write([]byte{statute.UserPassAuthVersion, statute.AuthFailure}); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return nil, fmt.Errorf("%s: zero user id returned", op)
	}

	if _, err := writer.Write([]byte{statute.UserPassAuthVersion, statute.AuthSuccess}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Done
	return &socks5.AuthContext{
		Method: statute.MethodUserPassAuth,
		Payload: map[string]string{
			"userId":   helper.Uint32ToString(userId),
			"username": username,
			"password": password,
		},
	}, nil
}
