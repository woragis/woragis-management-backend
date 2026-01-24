package csrf

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler handles CSRF token operations
type Handler struct {
	logger *slog.Logger
}

// NewHandler creates a new CSRF handler
func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func generateToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetCSRFToken returns CSRF token information. If the CSRF middleware is
// disabled or did not set a header, this handler will generate a token and
// return it (and also set it as a cookie and header) so API clients can
// obtain a usable token without middleware.
func (h *Handler) GetCSRFToken(c *fiber.Ctx) error {
	token := c.Get("X-CSRF-Token")
	if token == "" {
		// Generate a local token as middleware may be disabled in dev
		t, err := generateToken(32)
		if err != nil {
			h.logger.Error("failed to generate csrf token", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate csrf token",
			})
		}
		token = t

		// Set token in cookie so clients that read cookies can use it
		c.Cookie(&fiber.Cookie{
			Name:     "csrf_token",
			Value:    token,
			HTTPOnly: false,
			Secure:   false,
			SameSite: "Strict",
			MaxAge:   int((1 * time.Hour).Seconds()),
		})

		// Also set in response header for API clients
		c.Set("X-CSRF-Token", token)
	}

	// Return token in JSON body for ease of consumption by curl/dev clients
	return c.JSON(fiber.Map{
		"csrfToken": token,
		"message":  "CSRF token available in X-CSRF-Token header and csrf_token cookie",
	})
}