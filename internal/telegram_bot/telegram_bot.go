package telegrambot

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot              *tgbotapi.BotAPI
	onMessageHandler atomic.Value
}

func NewTelegramBot(apiKey string) (*TelegramBot, error) {
	config := config.GetConfig()
	bot, err := tgbotapi.NewBotAPI(config.TelegramConfig.APIKey)
	if err != nil {
		return nil, err
	}
	telegramBot := &TelegramBot{
		bot: bot,
	}
	return telegramBot, nil
}

func (b *TelegramBot) OnMessage(f func(update tgbotapi.Update) (reply tgbotapi.MessageConfig, err error)) {
	b.onMessageHandler.Store(f)
}

func (b *TelegramBot) onMessage(update tgbotapi.Update) (reply tgbotapi.MessageConfig, err error) {
	handler, ok := b.onMessageHandler.Load().(func(update tgbotapi.Update) (reply tgbotapi.MessageConfig, err error))
	if !ok {
		err = fmt.Errorf("Cannot cast handler")
		return
	}
	return handler(update)
}

func (b *TelegramBot) StartListen() error {
	b.bot.Debug = true

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	reply := make(chan tgbotapi.MessageConfig, 1024)

	go func() {
		for msg := range reply {
			b.bot.Send(msg)
		}
	}()

	go func() {
		for update := range b.bot.GetUpdatesChan(u) {
			if update.Message != nil { // If we got a message
				msg, err := b.onMessage(update)
				if err != nil {
					reply <- tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					continue
				}
				reply <- msg
			}
		}
	}()
	return nil
}
