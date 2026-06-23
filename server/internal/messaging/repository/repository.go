package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

type DestinationFilter struct {
	Channel    string
	ActiveOnly bool
	Query      string
}

func (r *Repository) ListDestinations(ctx context.Context, f DestinationFilter) ([]models.ChannelDestination, error) {
	var out []models.ChannelDestination
	q := r.db.WithContext(ctx).Order("name ASC")
	if f.ActiveOnly {
		q = q.Where("active = ?", true)
	}
	if f.Channel != "" {
		q = q.Where("channel = ?", f.Channel)
	}
	if term := strings.TrimSpace(f.Query); term != "" {
		like := "%" + term + "%"
		q = q.Where("name ILIKE ? OR external_id ILIKE ? OR description ILIKE ?", like, like, like)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list destinations: %w", err)
	}
	return out, nil
}

func (r *Repository) FindDestination(ctx context.Context, id uuid.UUID) (*models.ChannelDestination, error) {
	var row models.ChannelDestination
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find destination: %w", err)
	}
	return &row, nil
}

func (r *Repository) FindDestinationByExternal(ctx context.Context, channel, externalID string) (*models.ChannelDestination, error) {
	var row models.ChannelDestination
	err := r.db.WithContext(ctx).
		Where("channel = ? AND external_id = ?", channel, externalID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find destination by external: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateDestination(ctx context.Context, row *models.ChannelDestination) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	return nil
}

func (r *Repository) SaveDestination(ctx context.Context, row *models.ChannelDestination) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save destination: %w", err)
	}
	return nil
}

type TemplateFilter struct {
	DestinationID *uuid.UUID
	ProgramSlug     string
	ActiveOnly      bool
}

func (r *Repository) ListTemplates(ctx context.Context, f TemplateFilter) ([]models.MessageTemplate, error) {
	var out []models.MessageTemplate
	q := r.db.WithContext(ctx).Order("program_slug ASC, slug ASC")
	if f.ActiveOnly {
		q = q.Where("active = ?", true)
	}
	if f.DestinationID != nil {
		q = q.Where("destination_id = ?", *f.DestinationID)
	}
	if f.ProgramSlug != "" {
		q = q.Where("program_slug = ?", f.ProgramSlug)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	return out, nil
}

func (r *Repository) FindTemplate(ctx context.Context, id uuid.UUID) (*models.MessageTemplate, error) {
	var row models.MessageTemplate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find template: %w", err)
	}
	return &row, nil
}

func (r *Repository) FindTemplateBySlug(ctx context.Context, programSlug, slug string, destinationID *uuid.UUID) (*models.MessageTemplate, error) {
	var row models.MessageTemplate
	q := r.db.WithContext(ctx).Where("program_slug = ? AND slug = ?", programSlug, slug)
	if destinationID != nil {
		q = q.Where("destination_id = ?", *destinationID)
	} else {
		q = q.Where("destination_id IS NULL")
	}
	err := q.First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find template by slug: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateTemplate(ctx context.Context, row *models.MessageTemplate) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	return nil
}

func (r *Repository) SaveTemplate(ctx context.Context, row *models.MessageTemplate) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save template: %w", err)
	}
	return nil
}

func (r *Repository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.MessageTemplate{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete template: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) ListJobs(ctx context.Context, enabledOnly bool) ([]models.ScheduledJob, error) {
	var out []models.ScheduledJob
	q := r.db.WithContext(ctx).Order("name ASC")
	if enabledOnly {
		q = q.Where("enabled = ?", true)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	return out, nil
}

func (r *Repository) ListDueJobs(ctx context.Context, before time.Time) ([]models.ScheduledJob, error) {
	var out []models.ScheduledJob
	err := r.db.WithContext(ctx).
		Where("enabled = ? AND next_run_at IS NOT NULL AND next_run_at <= ?", true, before).
		Order("next_run_at ASC").
		Find(&out).Error
	if err != nil {
		return nil, fmt.Errorf("list due jobs: %w", err)
	}
	return out, nil
}

func (r *Repository) FindJob(ctx context.Context, id uuid.UUID) (*models.ScheduledJob, error) {
	var row models.ScheduledJob
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find job: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateJob(ctx context.Context, row *models.ScheduledJob) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create job: %w", err)
	}
	return nil
}

func (r *Repository) SaveJob(ctx context.Context, row *models.ScheduledJob) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save job: %w", err)
	}
	return nil
}

func (r *Repository) DeleteJob(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.ScheduledJob{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete job: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) CreateDelivery(ctx context.Context, row *models.MessageDelivery) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if row.SentAt.IsZero() {
		row.SentAt = time.Now().UTC()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create delivery: %w", err)
	}
	return nil
}

func (r *Repository) ListDeliveries(ctx context.Context, destinationID *uuid.UUID, limit int) ([]models.MessageDelivery, error) {
	var out []models.MessageDelivery
	q := r.db.WithContext(ctx).Order("sent_at DESC")
	if destinationID != nil {
		q = q.Where("destination_id = ?", *destinationID)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list deliveries: %w", err)
	}
	return out, nil
}
