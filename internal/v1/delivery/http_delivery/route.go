package httpdelivery

import (
	"net/http"
	"time"

	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.elastic.co/apm/module/apmechov4/v2"
)

func (h *HTTPDeliveryV1) initRoute() {
	v1 := h.echo.Group("/api/v1")
	for groupPath := range h.groupMapping {
		controllerGroup := v1.Group(groupPath)
		h.groupMapping[groupPath].InitControllerGroup(controllerGroup)
	}
}

func (h *HTTPDeliveryV1) initDefaulltMiddleware() {
	h.echo.Use(echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

	h.echo.OPTIONS("/*", func(c echo.Context) error { return nil })

	h.echo.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowCredentials: true,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		MaxAge: int((4 * time.Hour).Seconds()),

		UnsafeWildcardOriginWithAllowCredentials: true,
	}))

	h.echo.Use(apmechov4.Middleware())

	h.echo.Use(middleware.CompressWithConfig(middleware.CompressConfig{
		Level:       middleware.CompressLevelDefault,
		HandleError: true,
	}))

	h.echo.Use(echoMiddleware.CSRF())

	h.echo.Use(echoMiddleware.Decompress())

	h.echo.Use(middleware.LogRequest())
}
