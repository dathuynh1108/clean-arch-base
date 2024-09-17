package controller

import (
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"github.com/labstack/echo/v4"
)

type ErrorController struct {
	controller
}

func NewErrorController() *ErrorController {
	return &ErrorController{}
}

func (c *ErrorController) ErrorHandler(err error, ctx echo.Context) {
	if !ctx.Response().Committed {
		resErr := c.Failure(ctx, err)
		if resErr != nil {
			logger.GetLogger().Errorf("Error while responsing error: %v", resErr)
		}
	} else {
		logger.GetLogger().Infof("Request %s %s resulted in error: %v", ctx.Request().Method, ctx.Request().URL.Path, err)
	}
}
