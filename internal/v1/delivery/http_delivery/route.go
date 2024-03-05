package httpdelivery

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/controller"
	"github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func (h *httpDelivery) initRoute() {
	v1 := h.app.Group("/api/v1")
	for groupPath := range h.groupMapping {
		controllerGroup := v1.Group(groupPath)
		h.groupMapping[groupPath].InitControllerGroup(controllerGroup)
	}
}

func (h *httpDelivery) initDefaulltMiddleware() {
	h.app.Use(middleware.LogRequest)
	h.app.Use(recover.New())
}

func (h *httpDelivery) initWebSocket() {
	wsController := controller.ProvideWSController()
	h.app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	h.app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// Websocket logic
		wsController.Handle(c)
	}))
}

func (h *httpDelivery) initMetrics() {
	h.app.Get("/metrics", func(ctx *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(ctx.Context())
		return nil
	})
}
