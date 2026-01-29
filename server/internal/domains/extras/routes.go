package extras

import (
    "log/slog"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

// SetupRoutes registers extras routes under the provided router.
func SetupRoutes(api fiber.Router, db *gorm.DB, logger *slog.Logger) {
    repo := NewGormRepository(db)
    svc := NewService(repo, logger)
    handler := NewHandler(svc, logger)

    api.Post("/", handler.CreateExtra)
    api.Get("/", handler.ListExtras)
    api.Get("/:id", handler.GetExtra)
    api.Put("/:id", handler.UpdateExtra)
    api.Delete("/:id", handler.DeleteExtra)
}
