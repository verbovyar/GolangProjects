package botService

import (
	"fmt"
	"log"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type CmdHandler func(string) string

type Bot struct {
	bot   *botApi.BotAPI
	route map[string]CmdHandler
}

func Init(apiKey string) (*Bot, error) {
	bot, err := botApi.NewBotAPI(apiKey)
	if err != nil {
		return nil, errors.Wrap(err, "init tgbot")
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		bot:   bot,
		route: make(map[string]CmdHandler),
	}, nil
}

func (c *Bot) RegisterHandler(cmd string, f CmdHandler) {
	c.route[cmd] = f
}

func (c *Bot) Run() error {
	u := botApi.NewUpdate(0)
	u.Timeout = 60
	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := botApi.NewMessage(update.Message.Chat.ID, "")
		if cmd := update.Message.Command(); cmd != "" {
			if f, ok := c.route[cmd]; ok {
				msg.Text = f(update.Message.CommandArguments())
			} else {
				msg.Text = "Unknown command"
			}
		} else {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg.Text = fmt.Sprintf("you send <%v>", update.Message.Text)
		}
		_, err := c.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "send tg message")
		}
	}
	return nil
}
