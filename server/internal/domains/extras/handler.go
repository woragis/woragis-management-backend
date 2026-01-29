package extras

import (
    "log/slog"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"

    "woragis-management-service/pkg/middleware"
    "woragis-management-service/pkg/response"
)

// Handler exposes extras endpoints.
type Handler interface {
    CreateExtra(c *fiber.Ctx) error
    UpdateExtra(c *fiber.Ctx) error
    GetExtra(c *fiber.Ctx) error
    ListExtras(c *fiber.Ctx) error
    DeleteExtra(c *fiber.Ctx) error
}

type handler struct{
    service Service
    logger  *slog.Logger
}

// NewHandler constructs an extras handler.
func NewHandler(s Service, logger *slog.Logger) Handler {
    return &handler{service: s, logger: logger}
}

type createExtraPayload struct{
    Category string `json:"category,omitempty"`
    Text     string `json:"text,omitempty"`
    Ordinal  *int   `json:"ordinal,omitempty"`
}

type updateExtraPayload struct{
    Category *string `json:"category,omitempty"`
    Text     *string `json:"text,omitempty"`
    Ordinal  *int    `json:"ordinal,omitempty"`
}

func (h *handler) CreateExtra(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromFiberContext(c)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, fiber.Map{"message":"authentication required"})
    }

    var payload createExtraPayload
    if err := c.BodyParser(&payload); err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
    }

    extra, err := h.service.CreateExtra(c.Context(), userID, CreateExtraRequest(payload))
    if err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusCreated, extra)
}

func (h *handler) UpdateExtra(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromFiberContext(c)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, fiber.Map{"message":"authentication required"})
    }

    extraID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message":"invalid id"})
    }

    var payload updateExtraPayload
    if err := c.BodyParser(&payload); err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
    }

    extra, err := h.service.UpdateExtra(c.Context(), userID, extraID, UpdateExtraRequest(payload))
    if err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusOK, extra)
}

func (h *handler) GetExtra(c *fiber.Ctx) error {
    extraID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message":"invalid id"})
    }

    extra, err := h.service.GetExtra(c.Context(), extraID)
    if err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusOK, extra)
}

func (h *handler) ListExtras(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromFiberContext(c)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, fiber.Map{"message":"authentication required"})
    }

    list, err := h.service.ListExtras(c.Context(), userID)
    if err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusOK, list)
}

func (h *handler) DeleteExtra(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromFiberContext(c)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, fiber.Map{"message":"authentication required"})
    }

    extraID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message":"invalid id"})
    }

    if err := h.service.DeleteExtra(c.Context(), userID, extraID); err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusOK, fiber.Map{"message":"deleted"})
}

func (h *handler) handleError(c *fiber.Ctx, err error) error {
    if domainErr, ok := AsDomainError(err); ok {
        status := fiber.StatusInternalServerError
        switch domainErr.Code {
        case ErrCodeInvalidPayload:
            status = fiber.StatusBadRequest
        case ErrCodeNotFound:
            status = fiber.StatusNotFound
        case ErrCodeRepositoryFailure:
            status = fiber.StatusInternalServerError
        }
        return response.Error(c, status, domainErr.Code, fiber.Map{"message": domainErr.Message})
    }

    h.logger.Error("unexpected error in extras handler", slog.Any("error", err))
    return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{"message":"internal server error"})
}
