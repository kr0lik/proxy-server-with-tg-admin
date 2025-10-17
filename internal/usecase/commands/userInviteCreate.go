package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type userInviteCreate struct {
	storage StorageInterface
}

func (c *userInviteCreate) Id() string {
	return "invite"
}

func (c *userInviteCreate) IsForAdminOnly() bool { return true }

func (c *userInviteCreate) Arguments() []string {
	return []string{usernameArg}
}

func (c *userInviteCreate) Description() string {
	return "Create invite token for user"
}

func (c *userInviteCreate) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userInviteCreate.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	inviteToken := helper.GenerateInviteToken(username)

	if err := c.storage.UpdateInviteToken(username, inviteToken); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf("New token: `%s`\nExpired at %s", inviteToken, time.Now().Add(helper.InviteTokenTtl).Format(time.DateTime)), nil
}
