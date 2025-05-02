package commands

import (
	"errors"
	"fmt"
)

type StopUser struct {
	storage StorageInterface
}

func (c *StopUser) Id() string {
	return "deactivate"
}

func (c *StopUser) Arguments() string {
	return "{username}"
}

func (c *StopUser) Run(args ...string) (string, error) {
	var username string

	if len(args) == 0 {
		return "", errors.New("username is required")
	} else {
		username = args[0]
	}

	if err := c.storage.DeactivateUser(username); err != nil {
		return "", err
	}

	return fmt.Sprintf("User %s deactivated", username), nil
}
