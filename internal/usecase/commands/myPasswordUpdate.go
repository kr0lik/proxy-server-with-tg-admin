package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"proxy-server-with-tg-admin/internal/usecase/auth"
)

type myPasswordUpdate struct {
	ip            string
	port          uint
	storage       StorageInterface
	authenticator *auth.Authenticator
}

func (c *myPasswordUpdate) Id() string {
	return "passwd"
}

func (c *myPasswordUpdate) IsForAdminOnly() bool { return false }

func (c *myPasswordUpdate) Arguments() []string {
	return []string{"[password]"}
}

func (c *myPasswordUpdate) Description() string {
	return "Change my password"
}

func (c *myPasswordUpdate) Run(telegramId int64, args ...string) (string, error) {
	const op = "commands.myPasswordUpdate.Run"

	user, err := c.storage.GetUserByTelegramId(telegramId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !user.Active {
		return NotActiveAccountMsg, nil
	}

	password := helper.PasswordGenerate(len([]rune(user.Username)))

	if len(args) > 1 {
		password = args[1]
	}

	if err := c.storage.UpdatePassword(user.Username, password); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	c.authenticator.Forget(user.Username)

	return fmt.Sprintf("New credentials: \"*%s:%s@%s:%d*\"", user.Username, password, c.ip, c.port), nil
}
