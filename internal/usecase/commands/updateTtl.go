package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type UpdateTtl struct {
	storage StorageInterface
}

func (c *UpdateTtl) Id() string {
	return "ttl"
}

func (c *UpdateTtl) Arguments() []string {
	return []string{usernameArg, "[ttl]"}
}

func (c *UpdateTtl) Run(args ...string) (string, error) {
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
			return "", err
		}

		ttl = t
	}

	if err := c.storage.UpdateTtl(username, ttl); err != nil {
		return "", err
	}

	return fmt.Sprintf("User %s ttl updated to %s", username, helper.TtlToString(ttl)), nil
}
