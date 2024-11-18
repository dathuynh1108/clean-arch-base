package middleware

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"github.com/labstack/echo/v4"
)

func LogRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			var (
				startTime = time.Now()
			)

			defer func() {
				request := ctx.Request()

				entry := logger.GetLogger().
					WithFields(map[string]any{
						"method": request.Method,
						"path":   ctx.Path(),
						"ip":     ctx.RealIP(),
						"host":   request.Host,
						"time":   fmt.Sprintf("%d ms", time.Since(startTime).Milliseconds()),
					})

				if r := recover(); r != nil {
					var ok bool
					err, ok = r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					entry.WithField("stack", string(debug.Stack())).Error("Request panic recovered")
					return
				}

				if err != nil {
					entry.Error("Request failed")
					return
				}

				entry.Info("Request completed")
			}()

			err = next(ctx)
			return err
		}
	}
}
