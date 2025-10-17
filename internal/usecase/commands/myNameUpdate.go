package commands

import (
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type myNameUpdate struct {
	ip            string
	port          uint
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *myNameUpdate) Id() string {
	return "name"
}

func (c *myNameUpdate) IsForAdminOnly() bool { return false }

func (c *myNameUpdate) Arguments() []string {
	return []string{usernameArg}
}

func (c *myNameUpdate) Description() string {
	return "Change my username"
}

func (c *myNameUpdate) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userRename.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	user, err := c.storage.GetUserByTelegramId(telegramId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !user.Active {
		return NotActiveAccountMsg, nil
	}

	if err := c.storage.RenameUser(user.Username, username); err != nil {
		if errors.Is(err, ErrUserExists) {
			return ErrUserExists.Error(), nil
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.authenticator.Forget(user.Username)

	return fmt.Sprintf("New credentials: \"*%s:%s@%s:%d*\"", username, user.Password, c.ip, c.port), nil
}
