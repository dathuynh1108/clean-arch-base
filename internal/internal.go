package internal

import (
	"log"

	telegrambot "github.com/dathuynh1108/clean-arch-base/internal/telegram_bot"
	httpdelivery "github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartHTTPServer() {
	httpdelivery.ServeHTTP("0.0.0.0", "3000")
}

func StartTelegramBot() {
	config := config.GetConfig()
	telegramBot, err := telegrambot.NewTelegramBot(config.TelegramConfig.APIKey)
	if err != nil {
		panic(err)
	}

	telegramBot.OnMessage(func(update tgbotapi.Update) (reply tgbotapi.MessageConfig, err error) {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Reply for: "+update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		return msg, nil
	})

	telegramBot.StartListen()
}
