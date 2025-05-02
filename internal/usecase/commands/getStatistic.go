package commands

import (
	"errors"
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type GetStatistic struct {
	storage StorageInterface
}

func (c *GetStatistic) Id() string {
	return "stat"
}

func (c *GetStatistic) Arguments() string {
	return "{username}"
}

func (c *GetStatistic) Run(args ...string) (string, error) {
	var username string

	if len(args) == 0 {
		return "", errors.New("username is required")
	} else {
		username = args[0]
	}

	userStat, err := c.storage.GetStatistic(username)
	if err != nil {
		return "", err
	}

	updated := "-"
	if userStat.Updated.Unix() > 0 {
		updated = userStat.Updated.Format(time.DateOnly)
	}

	return fmt.Sprintf("Last activity %s.\nIn: %s\nOut: %s\nActive dyes: %d\n\nTotal In: %s\nTotal Out: %s",
		updated,
		helper.FromBytesFormat(userStat.TrafficInDay),
		helper.FromBytesFormat(userStat.TrafficOutDay),
		userStat.DaysActive,
		helper.FromBytesFormat(userStat.TrafficInTotal),
		helper.FromBytesFormat(userStat.TrafficOutTotal),
	), nil
}
