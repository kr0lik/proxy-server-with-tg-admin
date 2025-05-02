package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
)

type ListUsers struct {
	storage StorageInterface
}

func (c *ListUsers) Id() string {
	return "users"
}

func (c *ListUsers) Arguments() string {
	return ""
}

func (c *ListUsers) Run(args ...string) (string, error) {
	list := c.storage.ListUsers()

	res := ""

	for _, user := range list {
		active := "✅"

		if !user.Active {
			active = "⚫"
		}

		res += fmt.Sprintf("%s %s with ttl to %s\n", user.Username, active, helper.TtlToString(user.Ttl))
	}

	if res == "" {
		res = "Empty"
	}

	return res, nil
}
