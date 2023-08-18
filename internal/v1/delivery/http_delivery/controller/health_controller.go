package controller

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"
	"github.com/gofiber/fiber/v2"
)

type HealthControler struct {
	uc usecase.HealthUsecase
}

func NewHealthController(
	uc usecase.HealthUsecase,
) *HealthControler {
	return &HealthControler{
		uc: uc,
	}
}

func (h *HealthControler) GetHealth(ctx *fiber.Ctx) error {
	reply := h.uc.GetHealth()
	return ctx.SendString(reply)
}
