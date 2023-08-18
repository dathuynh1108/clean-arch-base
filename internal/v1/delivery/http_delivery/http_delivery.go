package httpdelivery

import (
	"fmt"
	"log"

	telegrambot "github.com/dathuynh1108/clean-arch-base/internal/telegram_bot"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/controller"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/repository"
	repositoryrouter "github.com/dathuynh1108/clean-arch-base/internal/v1/repository_router"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

type httpDelivery struct {
	app              *fiber.App
	healthController *controller.HealthControler
}

func StartHTTPServer(host, port string) error {
	config := config.GetConfig()

	// Repository
	healthRepo := repository.NewHealthRepository(nil)
	healthRepoRouter := repositoryrouter.NewRepositoryRouter[*repository.HealthRepository]()
	healthRepoRouter.SetRepo(repositoryrouter.AliasDefault, healthRepo)

	healthUC := usecase.NewHealthUsecase(healthRepoRouter)

	healthController := controller.NewHealthController(healthUC)
	httpDelivery := httpDelivery{
		app:              fiber.New(),
		healthController: healthController,
	}
	httpDelivery.initRoute()

	telegramBot, err := telegrambot.NewTelegramBot(config.TelegramConfig.APIKey)
	if err != nil {
		return err
	}

	telegramBot.OnMessage(func(update tgbotapi.Update) (reply tgbotapi.MessageConfig, err error) {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Reply for: "+update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		return msg, nil
	})
	telegramBot.StartListen()

	return httpDelivery.app.Listen(fmt.Sprintf("%v:%v", config.ServerConfig.Host, config.ServerConfig.Port))
}
