package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type createUser struct {
	ip      string
	port    uint
	storage StorageInterface
}

func (c *createUser) Id() string {
	return "create"
}

func (c *createUser) Arguments() []string {
	return []string{usernameArg, "[password]", "[ttl]"}
}

func (c *createUser) Run(args ...string) (string, error) {
	const op = "commands.createUser.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	password := helper.PasswordGenerate(len([]byte(username)))
	ttl := time.Now().Add(entity.DefaultUserTtl)

	secondInput := ""

	if len(args) > 2 {
		password = args[1]
		secondInput = args[2]
	} else if len(args) > 1 {
		secondInput = args[1]
	}

	if secondInput != "" {
		t, err := helper.StringToTtl(secondInput)
		if err != nil {
			password = secondInput
		} else {
			ttl = t
		}
	}

	_, err := c.storage.CreateUser(username, password)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := c.storage.ActivateUser(username); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := c.storage.UpdateTtl(username, ttl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	withTtl := helper.TtlToString(ttl)
	if withTtl != "" {
		withTtl = " with ttl to " + withTtl
	}

	return fmt.Sprintf("Created user: \"*%s:%s@%s:%d*\" %s\n", username, password, c.ip, c.port, withTtl), nil
}
