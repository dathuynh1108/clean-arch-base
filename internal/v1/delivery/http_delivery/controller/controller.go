package controller

import (
	"fmt"
	"net/http"

	"github.com/dathuynh1108/clean-arch-base/internal/v1/entity"
	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"
	"github.com/dathuynh1108/clean-arch-base/pkg/validator"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	BindAndValidate(ctx *fiber.Ctx, data any) error
	OK(ctx *fiber.Ctx, code int, message any, data any) error
	OKEmpty(ctx *fiber.Ctx) error
	Failure(ctx *fiber.Ctx, err error) error
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

func (c *controller) Failure(ctx *fiber.Ctx, err error) error {
	httpCode, code, message, errors := errorToResponse(err)
	return ctx.
		Status(httpCode).
		JSON(&entity.Response{
			Code:    code,
			Message: message,
			Data:    nil,
			Errors:  errors,
		})
}

func (c *controller) InitControllerGroup(app fiber.Router) {
	panic("InitControllerGroup is not implemented")
}

func errorToResponse(rootErr error) (httpCode int, code int, message any, errMessages []string) {
	switch errT := rootErr.(type) {
	case govalidator.ValidationErrors:
		httpCode = http.StatusBadRequest
		code = http.StatusBadRequest
		message = "Validation Error"
		errMessages = make([]string, len(errT))
		for i, fieldError := range errT {
			errMessages[i] = fmt.Sprintf(
				"Field validation for '%s' failed on the '%s' tag with value '%v'.",
				fieldError.Field(), fieldError.Tag(), fieldError.Value(),
			)
		}
		return
	default:
		httpCode = http.StatusInternalServerError
		code = http.StatusInternalServerError
		message = "Internal Server Error"
		errMessages = []string{rootErr.Error()}
		return
	}
}
