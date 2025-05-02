package commands

import (
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type ActivateUser struct {
	storage StorageInterface
}

func (c *ActivateUser) Id() string {
	return "activate"
}

func (c *ActivateUser) Arguments() string {
	return "{username} [ttl]"
}

func (c *ActivateUser) Run(args ...string) (string, error) {
	var username string

	if len(args) == 0 {
		return "", errors.New("username is required")
	} else {
		username = args[0]
	}

	ttl := time.Now().Add(entity.DefaultUserTtl)

	if len(args) > 1 {
		t, err := helper.StringToTtl(args[1])
		if err != nil {
			return "", err
		}
		ttl = t
	}

	if err := c.storage.ActivateUser(username); err != nil {
		return "", err
	}

	if err := c.storage.UpdateTtl(username, ttl); err != nil {
		return "", err
	}

	return fmt.Sprintf("User %s activated with ttl to %s", username, helper.TtlToString(ttl)), nil
}
