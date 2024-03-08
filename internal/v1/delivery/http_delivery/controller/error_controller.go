package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorController struct {
	controller
}

func NewErrorController() *ErrorController {
	return &ErrorController{}
}

func (c *ErrorController) ErrorHandler(ctx *fiber.Ctx, err error) error {
	return c.Failure(ctx, http.StatusInternalServerError, http.StatusInternalServerError, "Internal Server Error", []error{err})
}
