package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

const defaultSlug = "default"

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDefault(ctx context.Context) (*models.Profile, error) {
	var p models.Profile
	err := r.db.WithContext(ctx).Where("slug = ?", defaultSlug).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

func (r *Repository) EnsureDefault(ctx context.Context) (*models.Profile, error) {
	p, err := r.GetDefault(ctx)
	if err == nil {
		return p, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	p = &models.Profile{
		ID:           uuid.New(),
		Slug:         defaultSlug,
		DisplayName:  "Your Name",
		Availability: "not_available",
		SocialLinks:  []byte("[]"),
	}
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return nil, fmt.Errorf("create default profile: %w", err)
	}
	return p, nil
}

func (r *Repository) Save(ctx context.Context, p *models.Profile) error {
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return fmt.Errorf("save profile: %w", err)
	}
	return nil
}
