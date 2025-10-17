package telegram

import (
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"proxy-server-with-tg-admin/internal/usecase/commands"
	"strconv"
	"strings"
)

func MakeBot(token string, adminId int64, commands *commands.List, logger *slog.Logger) (*tele.Bot, error) {
	conf := tele.Settings{
		Token: token,
	}

	bot, err := tele.NewBot(conf)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	handleCommandStart(bot, commands, adminId)
	handleCommands(bot, commands, logger, adminId)
	handleMessages(bot, logger, adminId)

	return bot, nil
}

func handleCommandStart(bot *tele.Bot, commands *commands.List, adminId int64) {
	bot.Handle("/start", func(c tele.Context) error {
		isAdmin := c.Message().Sender.ID == adminId

		res := "Supported commands:\n"

		for _, cmd := range commands.List() {
			if !isAdmin && cmd.IsForAdminOnly() {
				continue
			}

			res += fmt.Sprintf("/%s %s - %s \n", cmd.Id(), strings.Join(cmd.Arguments(), " "), cmd.Description())
		}

		if !isAdmin {
			res += "\n\nType message there to contact admin."
		}

		return c.Reply(res)
	})
}

func handleCommands(bot *tele.Bot, commands *commands.List, logger *slog.Logger, adminId int64) {
	for _, cmd := range commands.List() {
		handleCommand(bot, cmd, logger, adminId)
	}
}

func handleCommand(bot *tele.Bot, cmd commands.Cmd, logger *slog.Logger, adminId int64) {
	bot.Handle("/"+cmd.Id(), func(c tele.Context) error {
		idAdmin := c.Message().Sender.ID == adminId

		if !idAdmin && cmd.IsForAdminOnly() {
			return nil
		}

		payload := c.Message().Payload

		args := strings.Fields(payload)

		res, err := cmd.Run(c.Message().Sender.ID, args...)
		if err != nil {
			if err := c.Reply(err.Error()); err != nil {
				if idAdmin {
					return c.Reply(err.Error())
				} else {
					logger.Error(err.Error())

					return c.Reply("Error. Contact with admin.")
				}
			}
		}

		if res != "" {
			if err := c.Reply(res, tele.ModeMarkdown); err != nil {
				if idAdmin {
					return c.Reply(err.Error())
				} else {
					logger.Error(err.Error())

					return c.Reply("Error. Contact with admin.")
				}
			}
		}

		return nil
	})
}

func handleMessages(bot *tele.Bot, logger *slog.Logger, adminId int64) {
	var currentReplyUserID int64

	bot.Handle(tele.OnText, func(c tele.Context) error {
		msg := c.Message()
		sender := msg.Sender

		if sender.ID != adminId {
			text := fmt.Sprintf("📩 from @%s (%d):\n%s", sender.Username, sender.ID, msg.Text)
			menu := &tele.ReplyMarkup{}
			replyBtn := menu.Data("Reply", fmt.Sprintf("reply_%d_%s", sender.ID, sender.Username))
			menu.Inline(menu.Row(replyBtn))

			_, err := bot.Send(tele.ChatID(adminId), text, menu)
			if err != nil {
				logger.Error("forward message to admin failed", "err", err)
			}

			return nil
		}

		if currentReplyUserID != 0 {
			_, err := bot.Send(tele.ChatID(currentReplyUserID), "📩 from Admin:\n"+msg.Text)
			if err != nil {
				logger.Error("forward admin answer failed", "err", err)
			}

			currentReplyUserID = 0
		}

		return nil
	})

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		cb := c.Callback()
		if cb.Sender.ID != adminId {
			return nil
		}

		parts := strings.Split(cb.Data, "_")
		if len(parts) < 3 { //nolint: mnd
			return c.Respond()
		}

		userID, _ := strconv.ParseInt(parts[1], 10, 64)
		currentReplyUserID = userID

		userName := parts[2]

		_ = c.Respond(&tele.CallbackResponse{Text: "Type message to user " + userName})

		return nil
	})
}
