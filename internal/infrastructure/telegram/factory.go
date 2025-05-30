package telegram

import (
	"fmt"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
	"proxy-server-with-tg-admin/internal/usecase/commands"
	"strings"
)

func MakeBot(token string, adminId int64, commands *commands.List) (*tele.Bot, error) {
	conf := tele.Settings{
		Token: token,
	}

	bot, err := tele.NewBot(conf)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	bot.Use(middleware.Whitelist(adminId))

	registerCommands(bot, commands)

	return bot, nil
}

func registerCommands(bot *tele.Bot, commands *commands.List) {
	bot.Handle("/start", func(c tele.Context) error {
		return c.Reply(getSupportedCommands(commands))
	})

	for _, cmd := range commands.List() {
		bot.Handle("/"+cmd.Id(), func(c tele.Context) error {
			payload := c.Message().Payload
			args := strings.Fields(payload)

			res, err := cmd.Run(args...)
			if res != "" {
				if err := c.Reply(res, tele.ModeMarkdown); err != nil {
					return c.Reply(err.Error())
				}
			}

			if err != nil {
				if err := c.Reply(err.Error()); err != nil {
					return c.Reply(err.Error())
				}
			}

			return nil
		})

		bot.Handle(tele.OnText, func(c tele.Context) error {
			text := c.Text()

			if text[0] == '/' {
				return c.Send("Unexpected command. Type /start to see available commands.")
			}

			return nil
		})
	}
}

func getSupportedCommands(commands *commands.List) string {
	res := "Supported commands:\n"

	for _, cmd := range commands.List() {
		res += fmt.Sprintf("/%s %s\n", cmd.Id(), strings.Join(cmd.Arguments(), " "))
	}

	return res
}
