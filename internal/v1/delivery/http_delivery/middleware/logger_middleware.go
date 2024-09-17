package middleware

import (
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

			defer func() {
				logger.GetLogger().
					WithFields(map[string]any{
						"method": request.Method,
						"path":   ctx.Path(),
						"ip":     ctx.RealIP(),
						"host":   ctx.Request().Host,
						"time":   time.Since(startTime).Milliseconds(),
					}).
					WithError(err).
					Info("Request")
			}()

			err = next(ctx)
			return err

		}
	}
}
