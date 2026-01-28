package certifications

import (
    "context"
    "time"

    "github.com/google/uuid"
    "log/slog"
)

// Service handles certifications use-cases.
type Service interface {
    CreateCertification(ctx context.Context, userID uuid.UUID, req CreateCertificationRequest) (*Certification, error)
    UpdateCertification(ctx context.Context, userID, certID uuid.UUID, req UpdateCertificationRequest) (*Certification, error)
    GetCertification(ctx context.Context, certID uuid.UUID) (*Certification, error)
    ListCertifications(ctx context.Context, userID uuid.UUID) ([]Certification, error)
    DeleteCertification(ctx context.Context, userID, certID uuid.UUID) error
}

type service struct{
    repo   Repository
    logger *slog.Logger
}

// NewService constructs a certifications service.
func NewService(repo Repository, logger *slog.Logger) Service {
    return &service{repo: repo, logger: logger}
}

type CreateCertificationRequest struct {
    Name        string `json:"name"`
    Issuer      string `json:"issuer,omitempty"`
    Date        string `json:"date,omitempty"`
    URL         string `json:"url,omitempty"`
    Description string `json:"description,omitempty"`
}

type UpdateCertificationRequest struct {
    Name        *string `json:"name,omitempty"`
    Issuer      *string `json:"issuer,omitempty"`
    Date        *string `json:"date,omitempty"`
    URL         *string `json:"url,omitempty"`
    Description *string `json:"description,omitempty"`
}

func (s *service) CreateCertification(ctx context.Context, userID uuid.UUID, req CreateCertificationRequest) (*Certification, error) {
    if userID == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
    }
    if req.Name == "" {
        return nil, NewDomainError(ErrCodeInvalidPayload, "certifications: name is required")
    }

    cert := NewCertification(userID, req.Name)
    if req.Issuer != "" {
        cert.Issuer = req.Issuer
    }
    if req.Date != "" {
        cert.Date = req.Date
    }
    if req.URL != "" {
        cert.URL = req.URL
    }
    if req.Description != "" {
        cert.Description = req.Description
    }

    if err := s.repo.CreateCertification(ctx, cert); err != nil {
        return nil, err
    }

    return cert, nil
}

func (s *service) UpdateCertification(ctx context.Context, userID, certID uuid.UUID, req UpdateCertificationRequest) (*Certification, error) {
    if userID == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
    }
    if certID == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }

    cert, err := s.repo.GetCertification(ctx, certID)
    if err != nil {
        return nil, err
    }
    if cert.UserID != userID {
        return nil, NewDomainError(ErrCodeInvalidPayload, "certifications: unauthorized")
    }

    if req.Name != nil {
        cert.Name = *req.Name
    }
    if req.Issuer != nil {
        cert.Issuer = *req.Issuer
    }
    if req.Date != nil {
        cert.Date = *req.Date
    }
    if req.URL != nil {
        cert.URL = *req.URL
    }
    if req.Description != nil {
        cert.Description = *req.Description
    }

    cert.UpdatedAt = time.Now().UTC()

    if err := s.repo.UpdateCertification(ctx, cert); err != nil {
        return nil, err
    }

    return cert, nil
}

func (s *service) GetCertification(ctx context.Context, certID uuid.UUID) (*Certification, error) {
    return s.repo.GetCertification(ctx, certID)
}

func (s *service) ListCertifications(ctx context.Context, userID uuid.UUID) ([]Certification, error) {
    return s.repo.ListCertificationsByUser(ctx, userID)
}

func (s *service) DeleteCertification(ctx context.Context, userID, certID uuid.UUID) error {
    if userID == uuid.Nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
    }
    if certID == uuid.Nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }

    cert, err := s.repo.GetCertification(ctx, certID)
    if err != nil {
        return err
    }
    if cert.UserID != userID {
        return NewDomainError(ErrCodeInvalidPayload, "certifications: unauthorized")
    }

    return s.repo.DeleteCertification(ctx, certID)
}
