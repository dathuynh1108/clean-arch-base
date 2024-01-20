package httpdelivery

import (
	"fmt"

	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/controller"
	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type httpDelivery struct {
	app              *fiber.App
	healthController *controller.HealthControler
	wsController     *controller.WSController
}

func ServeHTTP(host, port string) error {
	config := config.GetConfig()

	httpDelivery := httpDelivery{
		app: fiber.New(
			fiber.Config{
				ErrorHandler: controller.ProvideErrorController().ErrorHandler,
				JSONEncoder:  comjson.Marshal,
				JSONDecoder:  comjson.Unmarshal,
				Network:      fiber.NetworkTCP,
			},
		),
		healthController: controller.ProvideHealthController(),
		wsController:     controller.ProvideWSController(),
	}

	httpDelivery.initDefaulltMiddleware()
	httpDelivery.initRoute()

	return httpDelivery.app.Listen(fmt.Sprintf("%v:%v", config.ServerConfig.Host, config.ServerConfig.Port))
}
