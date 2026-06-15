package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]models.MediaAsset, error) {
	var out []models.MediaAsset
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list media: %w", err)
	}
	return out, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.MediaAsset, error) {
	var m models.MediaAsset
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find media: %w", err)
	}
	return &m, nil
}

func (r *Repository) Create(ctx context.Context, m *models.MediaAsset) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create media: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MediaAsset{})
	if res.Error != nil {
		return fmt.Errorf("delete media: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
