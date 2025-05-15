package commands

import (
	"fmt"
)

type clearStatistic struct {
	storage StorageInterface
}

func (c *clearStatistic) Id() string {
	return "clear"
}

func (c *clearStatistic) Arguments() []string {
	return []string{usernameArg}
}

func (c *clearStatistic) Run(args ...string) (string, error) {
	const op = "commands.clearStatistic.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	if err := c.storage.DeleteUserStat(username); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf("Statistic cleared for *%s*\n", username), nil
}
