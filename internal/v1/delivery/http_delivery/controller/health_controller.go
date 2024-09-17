package controller

import (
	"net/http"

	"github.com/dathuynh1108/clean-arch-base/internal/common"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"

	"github.com/labstack/echo/v4"
)

type HealthController struct {
	controller
	uc usecase.HealthUsecase
}

func NewHealthController(
	uc usecase.HealthUsecase,
) *HealthController {
	return &HealthController{
		uc: uc,
	}
}

func (h *HealthController) InitControllerGroup(app *echo.Group) {
	app.GET("*", h.GetHealth)
}

func (h *HealthController) GetHealth(c echo.Context) (err error) {
	var (
		ctx      = common.EchoWrapContext(c)
		reqModel = struct{}{}
	)

	if err = h.BindAndValidate(ctx, &reqModel); err != nil {
		return err
	}

	reply := h.uc.GetHealth(ctx)
	return h.OK(c, http.StatusOK, reply, reqModel)
}
