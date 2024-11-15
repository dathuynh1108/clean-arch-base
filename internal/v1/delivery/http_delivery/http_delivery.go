package httpdelivery

import (
	"context"

	"github.com/dathuynh1108/clean-arch-base/internal/common"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/controller"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"

	"github.com/labstack/echo/v4"
)

type HTTPDeliveryV1 struct {
	echo        *echo.Echo
	controllers []controller.Controller
}

func NewHTTPDeliveryV1() *HTTPDeliveryV1 {
	httpDelivery := HTTPDeliveryV1{
		echo: NewEchoDefault(),
		controllers: []controller.Controller{
			// Place other group here
			controller.ProvideHealthController(),
		},
	}

	return &httpDelivery
}
func (h *HTTPDeliveryV1) ServeHTTP(ctx context.Context, host string, port int) (stopper common.StartStopper, err error) {
	config := config.GetConfig()

	h.initDefaulltMiddleware()
	h.initRoute()

	echoSS := NewEchoStartStopper(h.echo, config.ServerConfig.Host, config.ServerConfig.Port)
	go func() {
		err = echoSS.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()

	return echoSS, nil
}
