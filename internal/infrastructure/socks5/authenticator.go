package socks5

import (
	"fmt"
	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
	"io"
	"proxy-server-with-tg-admin/internal/helper"
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

	// Verify the password
	userId := a.credentials.GetUserId(username, password, userAddr)

	if userId == 0 {
		if _, err := writer.Write([]byte{statute.UserPassAuthVersion, statute.AuthFailure}); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return nil, fmt.Errorf("%s: invalid username or password (%s:%s %s)", op, username, username, userAddr)
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
