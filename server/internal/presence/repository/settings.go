package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

const presenceSettingsID = "00000000-0000-0000-0000-000000000001"

func (r *Repository) GetSettings(ctx context.Context) (*models.PresenceSettings, error) {
	var row models.PresenceSettings
	err := r.db.WithContext(ctx).Where("id = ?", presenceSettingsID).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get presence settings: %w", err)
	}
	return &row, nil
}

func (r *Repository) EnsureSettings(ctx context.Context) (*models.PresenceSettings, error) {
	row, err := r.GetSettings(ctx)
	if err == nil {
		return row, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	id, _ := uuid.Parse(presenceSettingsID)
	row = &models.PresenceSettings{
		ID:               id,
		RemindersEnabled: true,
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, fmt.Errorf("create presence settings: %w", err)
	}
	return row, nil
}

func (r *Repository) SaveSettings(ctx context.Context, row *models.PresenceSettings) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save presence settings: %w", err)
	}
	return nil
}

func (r *Repository) ListDueReminders(ctx context.Context, before time.Time, limit int) ([]models.SocialPost, error) {
	if limit <= 0 {
		limit = 50
	}
	var out []models.SocialPost
	err := r.db.WithContext(ctx).
		Where("status = ?", models.SocialPostStatusScheduled).
		Where("scheduled_at IS NOT NULL AND scheduled_at <= ?", before).
		Where("reminder_sent_at IS NULL").
		Order("scheduled_at ASC").
		Limit(limit).
		Find(&out).Error
	if err != nil {
		return nil, fmt.Errorf("list due reminders: %w", err)
	}
	return out, nil
}
