package controller

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorController struct {
	controller
}

func NewErrorController() *ErrorController {
	return &ErrorController{}
}

func (c *ErrorController) ErrorHandler(ctx *fiber.Ctx, err error) error {
	return c.Failure(ctx, err)
}
