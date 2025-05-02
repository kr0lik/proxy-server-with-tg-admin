package commands

import (
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
)

type UpdatePassword struct {
	storage StorageInterface
}

func (c *UpdatePassword) Id() string {
	return "password"
}

func (c *UpdatePassword) Arguments() string {
	return "{username} [password]"
}

func (c *UpdatePassword) Run(args ...string) (string, error) {
	var username string

	if len(args) == 0 {
		return "", errors.New("username is required")
	} else {
		username = args[0]
	}

	password := helper.PasswordGenerate(len(username))

	if len(args) > 1 {
		password = args[1]
	}

	if err := c.storage.UpdatePassword(username, password); err != nil {
		return "", err
	}

	return fmt.Sprintf("New credentials: \"%s:%s\"", username, password), nil
}
