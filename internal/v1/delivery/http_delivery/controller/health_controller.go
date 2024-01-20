package controller

import (
	"net/http"

	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"
	"github.com/gofiber/fiber/v2"
)

type HealthControler struct {
	controller
	uc usecase.HealthUsecase
}

func NewHealthController(
	uc usecase.HealthUsecase,
) *HealthControler {
	return &HealthControler{
		uc: uc,
	}
}

func (h *HealthControler) InitControllerGroup(app fiber.Router) {
	app.Get("/health/get", h.GetHealth)
}

func (h *HealthControler) GetHealth(ctx *fiber.Ctx) error {
	reply := h.uc.GetHealth()
	return h.OK(ctx, http.StatusOK, "OK", reply)
}
