package middleware

import "github.com/gofiber/fiber/v2"

func LogRequest(ctx *fiber.Ctx) error {
	return ctx.Next()
}
