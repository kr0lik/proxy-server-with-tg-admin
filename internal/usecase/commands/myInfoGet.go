package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type myInfoGet struct {
	ip      string
	port    uint
	storage StorageInterface
}

func (c *myInfoGet) Id() string {
	return "info"
}

func (c *myInfoGet) IsForAdminOnly() bool { return false }

func (c *myInfoGet) Arguments() []string {
	return []string{}
}

func (c *myInfoGet) Description() string {
	return "Get my info"
}

func (c *myInfoGet) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.myInfoGet.Run"

	user, err := c.storage.GetUserByTelegramId(telegramId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !user.Active {
		return NotActiveAccountMsg, nil
	}

	userStat, err := c.storage.GetStatistic(user.ID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	updated := "-"
	if userStat.Updated.Unix() > 0 {
		updated = userStat.Updated.Format(time.DateOnly)
	}

	ttl := helper.TtlToString(user.Ttl)
	if ttl != "" {
		ttl = fmt.Sprintf("Time to live: %s\n", ttl)
	}

	res := fmt.Sprintf("*%s:%s@%s:%d*\n", user.Username, user.Password, c.ip, c.port) +
		ttl + "\n" +
		fmt.Sprintf("Last activity at %s.\n", updated) +
		fmt.Sprintf("in: %s\n", helper.BytesFormat(userStat.TrafficInDay)) +
		fmt.Sprintf("out: %s\n", helper.BytesFormat(userStat.TrafficOutDay)) +
		fmt.Sprintf("Active dyes: %d\n\n", userStat.DaysActive) +
		fmt.Sprintf("Total in: %s\n", helper.BytesFormat(userStat.TrafficInTotal)) +
		fmt.Sprintf("Total out: %s\n", helper.BytesFormat(userStat.TrafficOutTotal))

	return res, nil
}
