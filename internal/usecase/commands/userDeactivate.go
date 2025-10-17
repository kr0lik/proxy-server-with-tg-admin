package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type userDeactivate struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *userDeactivate) Id() string {
	return "deactivate"
}

func (c *userDeactivate) IsForAdminOnly() bool { return true }

func (c *userDeactivate) Arguments() []string {
	return []string{usernameArg}
}

func (c *userDeactivate) Description() string {
	return "Deactivate user"
}

func (c *userDeactivate) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userDeactivate.Run"
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
