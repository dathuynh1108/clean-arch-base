package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorControler struct {
	controller
}

func NewErrorController() *ErrorControler {
	return &ErrorControler{}
}

func (c *ErrorControler) ErrorHandler(ctx *fiber.Ctx, err error) error {
	return c.Failure(ctx, http.StatusInternalServerError, http.StatusInternalServerError, "Internal Server Error", []error{err})
}
