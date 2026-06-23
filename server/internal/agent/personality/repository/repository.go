package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Find(ctx context.Context) (*models.AgentPersonality, error) {
	var row models.AgentPersonality
	err := r.db.WithContext(ctx).Where("id = ?", models.DefaultPersonalityID).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find agent personality: %w", err)
	}
	return &row, nil
}

func (r *Repository) Save(ctx context.Context, row *models.AgentPersonality) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save agent personality: %w", err)
	}
	return nil
}

func (r *Repository) EnsureDefault(ctx context.Context) (*models.AgentPersonality, error) {
	row, err := r.Find(ctx)
	if err == nil {
		return row, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	def := models.DefaultAgentPersonality()
	if err := r.db.WithContext(ctx).Create(&def).Error; err != nil {
		return nil, fmt.Errorf("create default personality: %w", err)
	}
	return &def, nil
}
