package commands

import (
	"fmt"
)

type DeleteUser struct {
	storage StorageInterface
}

func (c *DeleteUser) Id() string {
	return "delete"
}

func (c *DeleteUser) Arguments() string {
	return "{username}"
}

func (c *DeleteUser) Run(args ...string) (string, error) {
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	if err := c.storage.DeleteUser(username); err != nil {
		return "", err
	}
	if err := c.storage.DeleteUserStat(username); err != nil {
		return "", err
	}

	return fmt.Sprintf("User %s deleted", username), nil
}
