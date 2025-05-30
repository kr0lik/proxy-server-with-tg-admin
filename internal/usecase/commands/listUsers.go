package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type listUsers struct {
	storage StorageInterface
}

func (c *listUsers) Id() string {
	return "users"
}

func (c *listUsers) Arguments() []string {
	return []string{}
}

func (c *listUsers) Run(args ...string) (string, error) {
	const op = "commands.listUsers.Run"

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

		res += fmt.Sprintf("%s *%s* %s\n", active, dto.Username, withTtl)
		res += fmt.Sprintf("Traffic in %s, out %s, dayes %d, last at %s\n", helper.BytesFormat(dto.TotalIn), helper.BytesFormat(dto.TotalOut), dto.DyesActive, dto.LastActive.Format(time.DateOnly))
	}

	if res == "" {
		res = "Empty"
	}

	if err != nil {
		return res, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
