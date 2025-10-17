package commands

import (
	"fmt"
)

type userStatisticClear struct {
	storage StorageInterface
}

func (c *userStatisticClear) Id() string {
	return "clear"
}

func (c *userStatisticClear) IsForAdminOnly() bool { return true }

func (c *userStatisticClear) Arguments() []string {
	return []string{usernameArg}
}

func (c *userStatisticClear) Description() string {
	return "Clear user statistic"
}

func (c *userStatisticClear) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userStatisticClear.Run"
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

	return fmt.Sprintf("Statistic cleared for *%s*\n", username), nil
}
