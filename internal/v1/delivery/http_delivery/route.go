package httpdelivery

import "github.com/gofiber/fiber/v2/middleware/recover"

func (h *httpDelivery) initRoute() {
	h.app.Get("/health/get", h.healthController.GetHealth)
}

func (h *httpDelivery) initDefaulltMiddleware() {
	h.app.Use(recover.New())
}
