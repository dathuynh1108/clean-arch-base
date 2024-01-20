package middleware

import (
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func LogRequest(ctx *fiber.Ctx) error {
	err := ctx.Next()
	if err != nil {
		logger.GetLogger().
			WithFields(map[string]any{
				"method": ctx.Method(),
				"path":   ctx.Path(),
				"ip":     ctx.IP(),
				"host":   ctx.Hostname(),
			}).
			WithError(err).
			Error("Request error")
	}
	return err
}
