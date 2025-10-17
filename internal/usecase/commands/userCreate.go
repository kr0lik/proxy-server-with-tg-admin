package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type userCreate struct {
	ip      string
	port    uint
	storage StorageInterface
}

func (c *userCreate) Id() string {
	return "create"
}

func (c *userCreate) IsForAdminOnly() bool { return true }

func (c *userCreate) Arguments() []string {
	return []string{usernameArg, "[password]", "[ttl]"}
}

func (c *userCreate) Description() string {
	return "Create user"
}

func (c *userCreate) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userCreate.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	password := helper.PasswordGenerate(len([]byte(username)))
	ttl := time.Now().Add(entity.DefaultUserTtl)

	secondInput := ""

	if len(args) > 2 { //nolint: mnd
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
