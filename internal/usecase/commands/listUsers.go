package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
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

	list, err := c.storage.ListUsers()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	res := ""

	for _, user := range list {
		active := "✅"

		if !user.Active {
			active = "⚫"
		}

		res += fmt.Sprintf("%s %s with ttl to %s\n", active, user.Username, helper.TtlToString(user.Ttl))
	}

	if res == "" {
		res = "Empty"
	}

	return res, nil
}
