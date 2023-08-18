package httpdelivery

func (h *httpDelivery) initRoute() error {
	h.app.Get("/health/get", h.healthController.GetHealth)
	return nil
}
