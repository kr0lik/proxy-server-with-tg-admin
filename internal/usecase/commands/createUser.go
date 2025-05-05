package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type CreateUser struct {
	storage StorageInterface
}

func (c *CreateUser) Id() string {
	return "create"
}

func (c *CreateUser) Arguments() string {
	return "{username} [password] [ttl]"
}

func (c *CreateUser) Run(args ...string) (string, error) {
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	password := helper.PasswordGenerate(len([]rune(username)))
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
		return "", err
	}

	if err := c.storage.ActivateUser(username); err != nil {
		return "", err
	}

	if err := c.storage.UpdateTtl(username, ttl); err != nil {
		return "", err
	}

	return fmt.Sprintf("Created user with credentials \"%s:%s\" and ttl to %s", username, password, helper.TtlToString(ttl)), nil
}
