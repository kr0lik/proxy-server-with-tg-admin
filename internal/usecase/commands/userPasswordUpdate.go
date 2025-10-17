package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type userPasswordUpdate struct {
	ip            string
	port          uint
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *userPasswordUpdate) Id() string {
	return "password"
}

func (c *userPasswordUpdate) IsForAdminOnly() bool { return true }

func (c *userPasswordUpdate) Arguments() []string {
	return []string{usernameArg, "[password]"}
}

func (c *userPasswordUpdate) Description() string {
	return "Change user password"
}

func (c *userPasswordUpdate) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userPasswordUpdate.Run"
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

	return fmt.Sprintf("New credentials: \"*%s:%s@%s:%d*\"", username, password, c.ip, c.port), nil
}
