package certifications

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes registers certifications routes under the provided router.
func SetupRoutes(api fiber.Router, db *gorm.DB, logger *slog.Logger) {
    repo := NewGormRepository(db)
    svc := NewService(repo, logger)
    handler := NewHandler(svc, logger)

    api.Post("/", handler.CreateCertification)
    api.Get("/", handler.ListCertifications)
    api.Get("/:id", handler.GetCertification)
    api.Put("/:id", handler.UpdateCertification)
    api.Delete("/:id", handler.DeleteCertification)
}
