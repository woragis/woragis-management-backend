package csrf

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers CSRF endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	api.Get("/", handler.GetCSRFToken)
}