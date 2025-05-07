package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type getStatistic struct {
	storage StorageInterface
}

func (c *getStatistic) Id() string {
	return "stat"
}

func (c *getStatistic) Arguments() []string {
	return []string{usernameArg}
}

func (c *getStatistic) Run(args ...string) (string, error) {
	const op = "commands.getStatistic.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	userStat, err := c.storage.GetStatistic(username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	updated := "-"
	if userStat.Updated.Unix() > 0 {
		updated = userStat.Updated.Format(time.DateOnly)
	}

	return fmt.Sprintf("Last activity %s.\nin: %s\nout: %s\nActive dyes: %d\n\nTotal in: %s\nTotal out: %s",
		updated,
		helper.FromBytesFormat(userStat.TrafficInDay),
		helper.FromBytesFormat(userStat.TrafficOutDay),
		userStat.DaysActive,
		helper.FromBytesFormat(userStat.TrafficInTotal),
		helper.FromBytesFormat(userStat.TrafficOutTotal),
	), nil
}
