package commands

import (
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
)

type join struct {
	storage StorageInterface
}

func (c *join) Id() string {
	return "join"
}

func (c *join) IsForAdminOnly() bool { return false }

func (c *join) Arguments() []string {
	return []string{"{token}"}
}

func (c *join) Description() string {
	return "Join by invite token"
}

func (c *join) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.join.Run"
	var token string

	if len(args) == 0 {
		return "Invite token required", nil
	} else {
		token = args[0]
	}

	if err := helper.CheckInviteToken(token); err != nil {
		if errors.Is(err, helper.ErrInviteTokenExpired) {
			return "Token expired. Request token again.", nil
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := c.storage.AssignTelegramIdByInviteToken(token, telegramId); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "You are joined!", nil
}
