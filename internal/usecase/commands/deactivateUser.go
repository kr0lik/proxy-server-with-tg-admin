package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type deactivateUser struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *deactivateUser) Id() string {
	return "deactivate"
}

func (c *deactivateUser) Arguments() []string {
	return []string{usernameArg}
}

func (c *deactivateUser) Run(args ...string) (string, error) {
	const op = "commands.deactivateUser.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	if err := c.storage.DeactivateUser(username); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.authenticator.Forget(username)

	return fmt.Sprintf("User *%s* deactivated", username), nil
}
