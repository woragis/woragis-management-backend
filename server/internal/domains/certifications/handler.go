package certifications

import (
    "log/slog"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"

    "woragis-management-service/pkg/middleware"
    "woragis-management-service/pkg/response"
)

// Handler exposes certifications endpoints.
type Handler interface {
    CreateCertification(c *fiber.Ctx) error
    UpdateCertification(c *fiber.Ctx) error
    GetCertification(c *fiber.Ctx) error
    ListCertifications(c *fiber.Ctx) error
    DeleteCertification(c *fiber.Ctx) error
}

type handler struct{
    service Service
    logger  *slog.Logger
}

// NewHandler constructs a certifications handler.
func NewHandler(s Service, logger *slog.Logger) Handler {
    return &handler{service: s, logger: logger}
}

type createCertificationPayload struct{
    Name string `json:"name"`
    Issuer string `json:"issuer,omitempty"`
    Date string `json:"date,omitempty"`
    URL string `json:"url,omitempty"`
    Description string `json:"description,omitempty"`
}

type updateCertificationPayload struct{
    Name *string `json:"name,omitempty"`
    Issuer *string `json:"issuer,omitempty"`
    Date *string `json:"date,omitempty"`
    URL *string `json:"url,omitempty"`
    Description *string `json:"description,omitempty"`
}

func (h *handler) CreateCertification(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromFiberContext(c)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, fiber.Map{"message":"authentication required"})
    }

    var payload createCertificationPayload
    if err := c.BodyParser(&payload); err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
    }

    cert, err := h.service.CreateCertification(c.Context(), userID, CreateCertificationRequest(payload))
    if err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusCreated, cert)
}

func (h *handler) UpdateCertification(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromFiberContext(c)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, fiber.Map{"message":"authentication required"})
    }

    certID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message":"invalid id"})
    }

    var payload updateCertificationPayload
    if err := c.BodyParser(&payload); err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
    }

    cert, err := h.service.UpdateCertification(c.Context(), userID, certID, UpdateCertificationRequest(payload))
    if err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusOK, cert)
}

func (h *handler) GetCertification(c *fiber.Ctx) error {
    certID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message":"invalid id"})
    }

    cert, err := h.service.GetCertification(c.Context(), certID)
    if err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusOK, cert)
}

func (h *handler) ListCertifications(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromFiberContext(c)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, fiber.Map{"message":"authentication required"})
    }

    list, err := h.service.ListCertifications(c.Context(), userID)
    if err != nil {
        return h.handleError(c, err)
    }

    return response.Success(c, fiber.StatusOK, list)
}

func (h *handler) DeleteCertification(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromFiberContext(c)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, fiber.Map{"message":"authentication required"})
    }

    certID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message":"invalid id"})
    }

    if err := h.service.DeleteCertification(c.Context(), userID, certID); err != nil {
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

    h.logger.Error("unexpected error in certifications handler", slog.Any("error", err))
    return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{"message":"internal server error"})
}
