package commands

import (
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

var ErrSecondUsername = errors.New("second username is required")

type renameUser struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *renameUser) Id() string {
	return "rename"
}

func (c *renameUser) Arguments() []string {
	return []string{usernameArg, usernameArg}
}

func (c *renameUser) Run(args ...string) (string, error) {
	const op = "commands.renameUser.Run"
	var username, usernameTo string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	if len(args) == 1 {
		return "", ErrSecondUsername
	} else {
		usernameTo = args[1]
	}

	if err := c.storage.RenameUser(username, usernameTo); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.authenticator.Forget(username)

	return fmt.Sprintf("User *%s* renamed to *%s*", username, usernameTo), nil
}
