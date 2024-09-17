package internal

import (
	"context"

	"github.com/dathuynh1108/clean-arch-base/internal/common"
	httpdelivery "github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/utils"
)

func StartHTTPServer(ctx context.Context) common.StartStopper {
	serverConfig := config.GetConfig().ServerConfig

	startStopper, err := httpdelivery.
		NewHTTPDeliveryV1().
		ServeHTTP(ctx, serverConfig.Host, serverConfig.Port)

	utils.PanicOnError(err)

	return startStopper
}
