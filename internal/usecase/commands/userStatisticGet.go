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

	res := fmt.Sprintf("*%s* statistic:\n", username) +
		fmt.Sprintf("Last activity at %s\n", updated) +
		fmt.Sprintf("Last day traffic in: %s\n", helper.BytesFormat(userStat.TrafficInDay)) +
		fmt.Sprintf("Last day traffic out: %s\n\n", helper.BytesFormat(userStat.TrafficOutDay)) +
		fmt.Sprintf("Active dayes: %d\n", userStat.DaysActive) +
		fmt.Sprintf("Total traffic in: %s\n", helper.BytesFormat(userStat.TrafficInTotal)) +
		fmt.Sprintf("Total traffic out: %s\n", helper.BytesFormat(userStat.TrafficOutTotal))

	return res, nil
}
