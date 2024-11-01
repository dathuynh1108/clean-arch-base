package middleware

import (
	"fmt"
	"time"

	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"github.com/labstack/echo/v4"
)

func LogRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			var (
				startTime = time.Now()
				request   = ctx.Request()
			)

			err = next(ctx)

			entry := logger.GetLogger().
				WithFields(map[string]any{
					"method": request.Method,
					"path":   ctx.Path(),
					"ip":     ctx.RealIP(),
					"host":   request.Host,
					"time":   fmt.Sprintf("%d ms", time.Since(startTime).Milliseconds()),
				}).
				WithError(err)

			if err != nil {
				entry.Error("Request failed")
			} else {
				entry.Info("Request completed")
			}

			return err
		}
	}
}
