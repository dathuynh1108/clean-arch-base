package controller

import (
	"fmt"
	"net/http"

	"github.com/dathuynh1108/clean-arch-base/internal/v1/entity"
	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller interface {
	BindAndValidate(ctx echo.Context, data any) error
	OK(ctx echo.Context, code int, message any, data any) error
	OKEmpty(ctx echo.Context) error
	Failure(ctx echo.Context, err error) error
	InitControllerGroup(app *echo.Group)
}

type controller struct{}

func (c *controller) BindAndValidate(ctx echo.Context, data any) error {
	if data == nil {
		return nil
	}
	if err := ctx.Bind(data); err != nil {
		return err
	}

	if err := ctx.Validate(data); err != nil {
		return err
	}

	return nil
}

func (c *controller) OK(ctx echo.Context, code int, message any, data any) error {
	return ctx.
		JSON(
			http.StatusOK,
			&entity.Response{
				Code:    code,
				Message: message,
				Data:    data,
			},
		)
}

func (c *controller) OKEmpty(ctx echo.Context) error {
	return ctx.
		JSON(
			http.StatusOK,
			&entity.Response{
				Code:    http.StatusOK,
				Message: "Success",
				Data:    nil,
			},
		)
}

func (c *controller) Failure(ctx echo.Context, err error) error {
	httpCode, code, message, errors := errorToResponse(err)
	return ctx.
		JSON(
			httpCode,
			&entity.Response{
				Code:    code,
				Message: message,
				Data:    nil,
				Errors:  errors,
			},
		)
}

func errorToResponse(err error) (httpCode int, code int, message any, errMessages []string) {
	switch errT := comerr.UnwrapFirst(err).(type) {
	case *echo.HTTPError:
		httpCode = errT.Code
		code = errT.Code
		message = errT.Message
		if errT.Internal != nil {
			errMessages = []string{errT.Internal.Error()}
		}
		return

	case govalidator.ValidationErrors:
		httpCode = http.StatusBadRequest
		code = http.StatusBadRequest
		message = "Validation Error"
		errMessages = make([]string, len(errT))
		for i, fieldError := range errT {
			errMessages[i] = fmt.Sprintf(
				"Field validation for '%s' failed on the '%s' tag with value '%v'",
				fieldError.Field(), fieldError.Tag(), fieldError.Value(),
			)
		}
		return
	default:
		httpCode = http.StatusBadRequest
		code = http.StatusBadRequest
		message = "Bad Request"
		errMessages = []string{err.Error()}
		return
	}
}
