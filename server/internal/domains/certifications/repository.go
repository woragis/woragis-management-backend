package certifications

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Repository defines persistence operations for certifications.
type Repository interface {
    CreateCertification(ctx context.Context, c *Certification) error
    UpdateCertification(ctx context.Context, c *Certification) error
    GetCertification(ctx context.Context, id uuid.UUID) (*Certification, error)
    ListCertificationsByUser(ctx context.Context, userID uuid.UUID) ([]Certification, error)
    DeleteCertification(ctx context.Context, id uuid.UUID) error
}

type gormRepository struct{
    db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository for certifications.
func NewGormRepository(db *gorm.DB) Repository {
    return &gormRepository{db: db}
}

func (r *gormRepository) CreateCertification(ctx context.Context, c *Certification) error {
    if c == nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrNilCertification)
    }
    now := time.Now().UTC()
    c.CreatedAt = now
    c.UpdatedAt = now
    if err := r.db.WithContext(ctx).Create(c).Error; err != nil {
        return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
    }
    return nil
}

func (r *gormRepository) UpdateCertification(ctx context.Context, c *Certification) error {
    if c == nil || c.ID == uuid.Nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }
    c.UpdatedAt = time.Now().UTC()
    res := r.db.WithContext(ctx).Model(&Certification{}).Where("id = ?", c.ID).Updates(c)
    if res.Error != nil {
        return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
    }
    if res.RowsAffected == 0 {
        return NewDomainError(ErrCodeNotFound, ErrNotFound)
    }
    return nil
}

func (r *gormRepository) GetCertification(ctx context.Context, id uuid.UUID) (*Certification, error) {
    if id == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }
    var c Certification
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, NewDomainError(ErrCodeNotFound, ErrNotFound)
        }
        return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
    }
    return &c, nil
}

func (r *gormRepository) ListCertificationsByUser(ctx context.Context, userID uuid.UUID) ([]Certification, error) {
    if userID == uuid.Nil {
        return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
    }
    var list []Certification
    if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error; err != nil {
        return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
    }
    return list, nil
}

func (r *gormRepository) DeleteCertification(ctx context.Context, id uuid.UUID) error {
    if id == uuid.Nil {
        return NewDomainError(ErrCodeInvalidPayload, ErrEmptyID)
    }
    res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Certification{})
    if res.Error != nil {
        return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
    }
    if res.RowsAffected == 0 {
        return NewDomainError(ErrCodeNotFound, ErrNotFound)
    }
    return nil
}
