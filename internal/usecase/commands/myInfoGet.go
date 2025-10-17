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

	res := fmt.Sprintf("*socks5://%s:%s@%s:%d*\n", user.Username, user.Password, c.ip, c.port) +
		ttl + "\n" +
		fmt.Sprintf("Last activity at %s\n", updated) +
		fmt.Sprintf("Last day traffic in: %s\n", helper.BytesFormat(userStat.TrafficInDay)) +
		fmt.Sprintf("Last day traffic out: %s\n\n", helper.BytesFormat(userStat.TrafficOutDay)) +
		fmt.Sprintf("Active dayes: %d\n", userStat.DaysActive) +
		fmt.Sprintf("Total traffic in: %s\n", helper.BytesFormat(userStat.TrafficInTotal)) +
		fmt.Sprintf("Total traffic out: %s\n", helper.BytesFormat(userStat.TrafficOutTotal))

	return res, nil
}
