package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type activateUser struct {
	storage StorageInterface
}

func (c *activateUser) Id() string {
	return "activate"
}

func (c *activateUser) Arguments() []string {
	return []string{usernameArg, "[ttl]"}
}

func (c *activateUser) Run(args ...string) (string, error) {
	const op = "commands.activateUser.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	ttl := time.Now().Add(entity.DefaultUserTtl)

	if len(args) > 1 {
		t, err := helper.StringToTtl(args[1])
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		ttl = t
	}

	if err := c.storage.ActivateUser(username); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := c.storage.UpdateTtl(username, ttl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf("User %s activated with ttl to %s", username, helper.TtlToString(ttl)), nil
}
