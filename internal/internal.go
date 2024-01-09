package internal

import (
	httpdelivery "github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
)

func StartHTTPServer() {
	serverConfig := config.GetConfig().ServerConfig
	err := httpdelivery.ServeHTTP(serverConfig.Host, serverConfig.Port)
	if err != nil {
		panic(err)
	}
}
