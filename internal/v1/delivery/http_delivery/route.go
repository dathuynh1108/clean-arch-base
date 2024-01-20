package httpdelivery

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/controller"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/middleware"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (h *httpDelivery) initRoute() {
	groupMapping := map[string]controller.Controller{
		"/health": h.healthController,
	}

	v1 := h.app.Group("/api/v1")
	for groupPath := range groupMapping {
		controllerGroup := v1.Group(groupPath)
		groupMapping[groupPath].InitControllerGroup(controllerGroup)
	}
}

func (h *httpDelivery) initDefaulltMiddleware() {
	h.app.Use(recover.New())
	h.app.Use(middleware.LogRequest)
}
