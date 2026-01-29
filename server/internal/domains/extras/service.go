package extras

import (
    "context"
    "time"

    "log/slog"

    "github.com/google/uuid"
)

// Service handles extras use-cases.
type Service interface {
    CreateExtra(ctx context.Context, userID uuid.UUID, req CreateExtraRequest) (*Extra, error)
    UpdateExtra(ctx context.Context, userID, extraID uuid.UUID, req UpdateExtraRequest) (*Extra, error)
    GetExtra(ctx context.Context, extraID uuid.UUID) (*Extra, error)
    ListExtras(ctx context.Context, userID uuid.UUID) ([]Extra, error)
    DeleteExtra(ctx context.Context, userID, extraID uuid.UUID) error
}

type service struct{
    repo   Repository
    logger *slog.Logger
}

// NewService constructs an extras service.
func NewService(repo Repository, logger *slog.Logger) Service {
    return &service{repo: repo, logger: logger}
}

type CreateExtraRequest struct {
    Category string `json:"category,omitempty"`
    Text     string `json:"text,omitempty"`
    Ordinal  *int   `json:"ordinal,omitempty"`
}

type UpdateExtraRequest struct {
    Category *string `json:"category,omitempty"`
    Text     *string `json:"text,omitempty"`
    Ordinal  *int    `json:"ordinal,omitempty"`
}

func (s *service) CreateExtra(ctx context.Context, userID uuid.UUID, req CreateExtraRequest) (*Extra, error) {
    if userID == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
    }

    e := NewExtra(userID, req.Category, req.Text)
    if req.Ordinal != nil {
        e.Ordinal = *req.Ordinal
    }

    if err := s.repo.CreateExtra(ctx, e); err != nil {
        return nil, err
    }
    return e, nil
}

func (s *service) UpdateExtra(ctx context.Context, userID, extraID uuid.UUID, req UpdateExtraRequest) (*Extra, error) {
    if userID == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
    }
    if extraID == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }

    e, err := s.repo.GetExtra(ctx, extraID)
    if err != nil {
        return nil, err
    }
    if e.UserID != userID {
        return nil, NewDomainError(ErrCodeInvalidPayload, "extras: unauthorized")
    }

    if req.Category != nil {
        e.Category = *req.Category
    }
    if req.Text != nil {
        e.Text = *req.Text
    }
    if req.Ordinal != nil {
        e.Ordinal = *req.Ordinal
    }

    e.UpdatedAt = time.Now().UTC()

    if err := s.repo.UpdateExtra(ctx, e); err != nil {
        return nil, err
    }
    return e, nil
}

func (s *service) GetExtra(ctx context.Context, extraID uuid.UUID) (*Extra, error) {
    return s.repo.GetExtra(ctx, extraID)
}

func (s *service) ListExtras(ctx context.Context, userID uuid.UUID) ([]Extra, error) {
    return s.repo.ListExtrasByUser(ctx, userID)
}

func (s *service) DeleteExtra(ctx context.Context, userID, extraID uuid.UUID) error {
    if userID == uuid.Nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
    }
    if extraID == uuid.Nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }

    e, err := s.repo.GetExtra(ctx, extraID)
    if err != nil {
        return err
    }
    if e.UserID != userID {
        return NewDomainError(ErrCodeInvalidPayload, "extras: unauthorized")
    }

    return s.repo.DeleteExtra(ctx, extraID)
}
