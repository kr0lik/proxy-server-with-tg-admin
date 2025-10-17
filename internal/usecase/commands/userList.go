package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type userList struct {
	storage StorageInterface
}

func (c *userList) Id() string {
	return "users"
}

func (c *userList) IsForAdminOnly() bool { return true }

func (c *userList) Arguments() []string {
	return []string{}
}

func (c *userList) Description() string {
	return "List users"
}

func (c *userList) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userList.Run"

	list, err := c.storage.ListUsersWithStat()

	res := ""

	for _, dto := range list {
		active := "✅"

		if !dto.Active {
			active = "⚫"
		}

		withTtl := helper.TtlToString(dto.Ttl)
		if withTtl != "" {
			withTtl = " with ttl to " + withTtl
		}

		lastAt := "-"
		if dto.LastActive.Unix() > 0 {
			lastAt = dto.LastActive.Format(time.DateOnly)
		}

		hasTg := ""
		if dto.TelegramId > 0 {
			hasTg = "tg"
		}

		res += fmt.Sprintf("%s *%s* %s %s\n", active, dto.Username, withTtl, hasTg)
		res += fmt.Sprintf("Traffic in %s, out %s, dayes %d, last at %s\n", helper.BytesFormat(dto.TotalIn), helper.BytesFormat(dto.TotalOut), dto.DyesActive, lastAt)
	}

	if res == "" {
		res = "Empty"
	}

	if err != nil {
		return res, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
