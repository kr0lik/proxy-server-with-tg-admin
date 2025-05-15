package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type deleteUser struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *deleteUser) Id() string {
	return "delete"
}

func (c *deleteUser) Arguments() []string {
	return []string{usernameArg}
}

func (c *deleteUser) Run(args ...string) (string, error) {
	const op = "commands.deleteUser.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	if err := c.storage.DeleteUserWithStat(username); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.authenticator.Forget(username)

	return fmt.Sprintf("User *%s* deleted", username), nil
}
