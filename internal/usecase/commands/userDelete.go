package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type userDelete struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *userDelete) Id() string {
	return "delete"
}

func (c *userDelete) IsForAdminOnly() bool { return true }

func (c *userDelete) Arguments() []string {
	return []string{usernameArg}
}

func (c *userDelete) Description() string {
	return "Delete user"
}

func (c *userDelete) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userDelete.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	userId, err := c.storage.GetUserIdByUsername(username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := c.storage.DeleteUserStat(userId); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := c.storage.DeleteUser(userId); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.authenticator.Forget(username)

	return fmt.Sprintf("User *%s* deleted", username), nil
}
