package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/helper"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"time"
)

type updateTtl struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *updateTtl) Id() string {
	return "ttl"
}

func (c *updateTtl) Arguments() []string {
	return []string{usernameArg, "[ttl]"}
}

func (c *updateTtl) Run(args ...string) (string, error) {
	const op = "commands.updateTtl.Run"
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

	if err := c.storage.UpdateTtl(username, ttl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.authenticator.UpdateUserTtl(username, ttl)

	return fmt.Sprintf("User %s ttl updated to %s", username, helper.TtlToString(ttl)), nil
}
