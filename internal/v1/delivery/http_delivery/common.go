package httpdelivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dathuynh1108/clean-arch-base/internal/common"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/controller"
	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"
	"github.com/dathuynh1108/clean-arch-base/pkg/validator"
	"github.com/labstack/echo/v4"
)

type EchoStartStopper struct {
	echo      *echo.Echo
	host      string
	port      int
	isRunning bool
}

func NewEchoStartStopper(e *echo.Echo, host string, port int) *EchoStartStopper {
	return &EchoStartStopper{
		echo:      e,
		host:      host,
		port:      port,
		isRunning: false,
	}
}

func (es *EchoStartStopper) Name() string {
	return "echo"
}

func (es *EchoStartStopper) Start(context.Context) error {
	es.isRunning = true
	defer func() { es.isRunning = false }()

	address := fmt.Sprintf("%s:%d", es.host, es.port)
	err := es.echo.Start(address)
	if err != nil {
		if err != http.ErrServerClosed {
			return comerr.
				WrapMessage(err, "echo: start server")
		}
	}
	return nil
}

func (es *EchoStartStopper) BeforeStop(context.Context) {
	es.isRunning = false
}

func (es *EchoStartStopper) Stop(ctx context.Context) error {
	return es.echo.Shutdown(ctx)
}

func NewEchoDefault() *echo.Echo {
	echoObj := echo.New()

	echoObj.HTTPErrorHandler = controller.ProvideErrorController().ErrorHandler
	echoObj.JSONSerializer = common.EchoJSONSerializer{}
	echoObj.Binder = &common.EchoBinder{}
	echoObj.Validator = validator.GetValidator()

	return echoObj
}
