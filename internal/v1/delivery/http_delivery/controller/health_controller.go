package controller

import (
	"net/http"

	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"
	"github.com/gofiber/fiber/v2"
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

func (h *HealthController) InitControllerGroup(app fiber.Router) {
	app.Get("/health/*", h.GetHealth)
}

func (h *HealthController) GetHealth(c *fiber.Ctx) (err error) {
	if err = h.BindAndValidate(c, nil); err != nil {
		return err
	}
	var (
		ctx = c.Context()
	)

	reply := h.uc.GetHealth(ctx)
	return h.OK(c, http.StatusOK, "OK", reply)
}
