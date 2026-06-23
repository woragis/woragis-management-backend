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

type ListFilter struct {
	Query        string
	Relationship string
	Organization string
	Stage        string
	ProjectID    *uuid.UUID
	ActiveOnly   bool
}

func (r *Repository) ListContacts(ctx context.Context, f ListFilter) ([]models.Contact, error) {
	var out []models.Contact
	q := r.db.WithContext(ctx).Order("name ASC, organization ASC")
	if f.ActiveOnly {
		q = q.Where("active = ?", true)
	}
	if f.Relationship != "" {
		q = q.Where("relationship = ?", f.Relationship)
	}
	if f.Organization != "" {
		q = q.Where("organization ILIKE ?", "%"+escapeLike(f.Organization)+"%")
	}
	if f.Stage != "" {
		q = q.Where("stage = ?", f.Stage)
	}
	if f.ProjectID != nil {
		q = q.Where("project_id = ?", *f.ProjectID)
	}
	if term := strings.TrimSpace(f.Query); term != "" {
		like := "%" + escapeLike(term) + "%"
		q = q.Where(
			"name ILIKE ? OR organization ILIKE ? OR role_title ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR display_name ILIKE ?",
			like, like, like, like, like, like,
		)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list contacts: %w", err)
	}
	return out, nil
}

func (r *Repository) FindContact(ctx context.Context, id uuid.UUID) (*models.Contact, error) {
	var row models.Contact
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find contact: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateContact(ctx context.Context, row *models.Contact) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create contact: %w", err)
	}
	return nil
}

func (r *Repository) SaveContact(ctx context.Context, row *models.Contact) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save contact: %w", err)
	}
	return nil
}

func (r *Repository) ListInteractions(ctx context.Context, contactID uuid.UUID) ([]models.ContactInteraction, error) {
	var out []models.ContactInteraction
	err := r.db.WithContext(ctx).
		Where("contact_id = ?", contactID).
		Order("happened_at DESC, created_at DESC").
		Find(&out).Error
	if err != nil {
		return nil, fmt.Errorf("list interactions: %w", err)
	}
	return out, nil
}

func (r *Repository) FindInteraction(ctx context.Context, contactID, id uuid.UUID) (*models.ContactInteraction, error) {
	var row models.ContactInteraction
	err := r.db.WithContext(ctx).
		Where("id = ? AND contact_id = ?", id, contactID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find interaction: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateInteraction(ctx context.Context, row *models.ContactInteraction) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create interaction: %w", err)
	}
	return nil
}

func (r *Repository) ListContactsDueFollowUp(ctx context.Context, before time.Time) ([]models.Contact, error) {
	var out []models.Contact
	err := r.db.WithContext(ctx).
		Where("active = ? AND next_follow_up_at IS NOT NULL AND next_follow_up_at <= ?", true, before).
		Order("next_follow_up_at ASC").
		Find(&out).Error
	if err != nil {
		return nil, fmt.Errorf("list contacts due follow-up: %w", err)
	}
	return out, nil
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
