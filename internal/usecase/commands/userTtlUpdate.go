package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/helper"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"time"
)

type userTtlUpdate struct {
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *userTtlUpdate) Id() string {
	return "ttl"
}

func (c *userTtlUpdate) IsForAdminOnly() bool { return true }

func (c *userTtlUpdate) Arguments() []string {
	return []string{usernameArg, "[ttl]"}
}

func (c *userTtlUpdate) Description() string {
	return "Change ttl for user"
}

func (c *userTtlUpdate) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userTtlUpdate.Run"
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

	withTtl := helper.TtlToString(ttl)
	if withTtl == "" {
		withTtl = "unlimited"
	}

	return fmt.Sprintf("InviteTokenTtl updated for *%s* to %s\n", username, withTtl), nil
}
