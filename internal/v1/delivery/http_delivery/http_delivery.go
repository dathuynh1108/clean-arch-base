package httpdelivery

import (
	"fmt"

	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/controller"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/repository"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"
	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/database"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/repo"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type httpDelivery struct {
	app              *fiber.App
	healthController *controller.HealthControler
}

func ServeHTTP(host, port string) error {
	config := config.GetConfig()

	// Repository
	healthRepoRouter := repo.NewRepoRouter(
		dbpool.DBDefault,
		database.ProvideDBPool(),
		func(db *gorm.DB) *repository.HealthRepository {
			return repository.NewHealthRepository(db)
		},
	)

	// Usecase
	healthUC := usecase.NewHealthUsecase(healthRepoRouter)

	// Controller
	errorController := controller.NewErrorController()
	healthController := controller.NewHealthController(healthUC)

	httpDelivery := httpDelivery{
		app: fiber.New(
			fiber.Config{
				ErrorHandler: errorController.ErrorHandler,
				JSONEncoder:  comjson.Marshal,
				JSONDecoder:  comjson.Unmarshal,
				Network:      fiber.NetworkTCP,
			},
		),
		healthController: healthController,
	}

	httpDelivery.initDefaulltMiddleware()
	httpDelivery.initRoute()

	return httpDelivery.app.Listen(fmt.Sprintf("%v:%v", config.ServerConfig.Host, config.ServerConfig.Port))
}
