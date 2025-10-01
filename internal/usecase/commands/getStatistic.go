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

	res := fmt.Sprintf("*%s* stats:\n", username) +
		fmt.Sprintf("Last activity at %s.\n", updated) +
		fmt.Sprintf("in: %s\n", helper.BytesFormat(userStat.TrafficInDay)) +
		fmt.Sprintf("out: %s\n", helper.BytesFormat(userStat.TrafficOutDay)) +
		fmt.Sprintf("Active dyes: %d\n\n", userStat.DaysActive) +
		fmt.Sprintf("Total in: %s\n", helper.BytesFormat(userStat.TrafficInTotal)) +
		fmt.Sprintf("Total out: %s\n", helper.BytesFormat(userStat.TrafficOutTotal))

	return res, nil
}
