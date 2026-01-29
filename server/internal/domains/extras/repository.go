package extras

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Repository defines persistence operations for extras.
type Repository interface {
    CreateExtra(ctx context.Context, e *Extra) error
    UpdateExtra(ctx context.Context, e *Extra) error
    GetExtra(ctx context.Context, id uuid.UUID) (*Extra, error)
    ListExtrasByUser(ctx context.Context, userID uuid.UUID) ([]Extra, error)
    DeleteExtra(ctx context.Context, id uuid.UUID) error
}

type gormRepository struct{
    db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository for extras.
func NewGormRepository(db *gorm.DB) Repository {
    return &gormRepository{db: db}
}

func (r *gormRepository) CreateExtra(ctx context.Context, e *Extra) error {
    if e == nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrNilExtra)
    }
    now := time.Now().UTC()
    e.CreatedAt = now
    e.UpdatedAt = now
    if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
        return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
    }
    return nil
}

func (r *gormRepository) UpdateExtra(ctx context.Context, e *Extra) error {
    if e == nil || e.ID == uuid.Nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }
    e.UpdatedAt = time.Now().UTC()
    res := r.db.WithContext(ctx).Model(&Extra{}).Where("id = ?", e.ID).Updates(e)
    if res.Error != nil {
        return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
    }
    if res.RowsAffected == 0 {
        return NewDomainError(ErrCodeNotFound, ErrNotFound)
    }
    return nil
}

func (r *gormRepository) GetExtra(ctx context.Context, id uuid.UUID) (*Extra, error) {
    if id == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }
    var e Extra
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&e).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, NewDomainError(ErrCodeNotFound, ErrNotFound)
        }
        return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
    }
    return &e, nil
}

func (r *gormRepository) ListExtrasByUser(ctx context.Context, userID uuid.UUID) ([]Extra, error) {
    if userID == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
    }
    var list []Extra
    if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("ordinal asc, created_at desc").Find(&list).Error; err != nil {
        return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
    }
    return list, nil
}

func (r *gormRepository) DeleteExtra(ctx context.Context, id uuid.UUID) error {
    if id == uuid.Nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }
    res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Extra{})
    if res.Error != nil {
        return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
    }
    if res.RowsAffected == 0 {
        return NewDomainError(ErrCodeNotFound, ErrNotFound)
    }
    return nil
}
