package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"time"
)

type userStatisticGet struct {
	storage StorageInterface
}

func (c *userStatisticGet) Id() string {
	return "stat"
}

func (c *userStatisticGet) IsForAdminOnly() bool { return true }

func (c *userStatisticGet) Arguments() []string {
	return []string{usernameArg}
}

func (c *userStatisticGet) Description() string {
	return "Get user statistic"
}

func (c *userStatisticGet) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.userStatisticGet.Run"
	var username string

	if len(args) == 0 {
		return "", ErrUsernameRequired
	} else {
		username = args[0]
	}

	userId, err := c.storage.GetUserIdByUsername(username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	userStat, err := c.storage.GetStatistic(userId)
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
