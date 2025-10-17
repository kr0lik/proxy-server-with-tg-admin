package commands

import (
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

var ErrSecondUsername = errors.New("second username is required")

type userRename struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *userRename) Id() string {
	return "rename"
}

func (c *userRename) IsForAdminOnly() bool { return true }

func (c *userRename) Arguments() []string {
	return []string{usernameArg, usernameArg}
}

func (c *userRename) Description() string {
	return "Change username"
}

func (c *userRename) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userRename.Run"
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
