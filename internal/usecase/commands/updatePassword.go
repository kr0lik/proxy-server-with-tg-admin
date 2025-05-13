package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type updatePassword struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *updatePassword) Id() string {
	return "password"
}

func (c *updatePassword) Arguments() []string {
	return []string{usernameArg, "[password]"}
}

func (c *updatePassword) Run(args ...string) (string, error) {
	const op = "commands.updatePassword.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	password := helper.PasswordGenerate(len([]rune(username)))

	if len(args) > 1 {
		password = args[1]
	}

	if err := c.storage.UpdatePassword(username, password); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.authenticator.Forget(username)

	return fmt.Sprintf("New credentials: \"%s:%s\"", username, password), nil
}
