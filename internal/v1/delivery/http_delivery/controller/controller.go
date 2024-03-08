package controller

import (
	"net/http"

	"github.com/dathuynh1108/clean-arch-base/internal/v1/entity"
	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"
	"github.com/dathuynh1108/clean-arch-base/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	BindAndValidate(ctx *fiber.Ctx, data any) error
	OK(ctx *fiber.Ctx, code int, message any, data any) error
	OKEmpty(ctx *fiber.Ctx) error
	Failure(ctx *fiber.Ctx, httpCode int, code int, message any, errors []error) error
	InitControllerGroup(app fiber.Router)
}

type controller struct{}

func (c *controller) BindAndValidate(ctx *fiber.Ctx, data any) error {
	if data == nil {
		return nil
	}

	if err := ctx.BodyParser(data); err != nil {
		return comerr.WrapError(err, "Failed to parse request body")
	}
	err := validator.GetValidator().Validate(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) OK(ctx *fiber.Ctx, code int, message any, data any) error {
	return ctx.
		Status(http.StatusOK).
		JSON(&entity.Response{
			Code:    code,
			Message: message,
			Data:    data,
		})
}

func (c *controller) OKEmpty(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(nil)
}

func (c *controller) Failure(ctx *fiber.Ctx, httpCode int, code int, message any, errors []error) error {
	errorsString := make([]string, len(errors))
	for i, err := range errors {
		errorsString[i] = err.Error()
	}

	return ctx.
		Status(httpCode).
		JSON(&entity.Response{
			Code:    code,
			Message: message,
			Data:    nil,
			Errors:  errorsString,
		})
}

func (c *controller) InitControllerGroup(app fiber.Router) {
	panic("InitControllerGroup is not implemented")
}
