package commands

import (
	"fmt"
)

type deleteUser struct {
	storage StorageInterface
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

	if err := c.storage.DeleteUser(username); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := c.storage.DeleteUserStat(username); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf("User %s deleted", username), nil
}
