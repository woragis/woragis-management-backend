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

type CampaignFilter struct {
	Goal       string
	ProjectID  *uuid.UUID
	ActiveOnly bool
}

type PostFilter struct {
	Platform   string
	Goal       string
	Status     string
	ProjectID  *uuid.UUID
	CampaignID *uuid.UUID
	Limit      int
}

type TemplateFilter struct {
	Platform   string
	Goal       string
	ActiveOnly bool
}

func (r *Repository) ListCampaigns(ctx context.Context, f CampaignFilter) ([]models.SocialCampaign, error) {
	var out []models.SocialCampaign
	q := r.db.WithContext(ctx).Order("created_at DESC")
	if f.ActiveOnly {
		q = q.Where("active = ?", true)
	}
	if f.Goal != "" {
		q = q.Where("goal = ?", f.Goal)
	}
	if f.ProjectID != nil {
		q = q.Where("project_id = ?", *f.ProjectID)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list campaigns: %w", err)
	}
	return out, nil
}

func (r *Repository) FindCampaign(ctx context.Context, id uuid.UUID) (*models.SocialCampaign, error) {
	var row models.SocialCampaign
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find campaign: %w", err)
	}
	return &row, nil
}

func (r *Repository) SaveCampaign(ctx context.Context, row *models.SocialCampaign) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save campaign: %w", err)
	}
	return nil
}

func (r *Repository) CreateCampaign(ctx context.Context, row *models.SocialCampaign) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	return r.SaveCampaign(ctx, row)
}

func (r *Repository) DeleteCampaign(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.SocialCampaign{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete campaign: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) ListTemplates(ctx context.Context, f TemplateFilter) ([]models.PostTemplate, error) {
	var out []models.PostTemplate
	q := r.db.WithContext(ctx).Order("name ASC")
	if f.ActiveOnly {
		q = q.Where("active = ?", true)
	}
	if f.Platform != "" {
		q = q.Where("platform IN ?", []string{f.Platform, "any"})
	}
	if f.Goal != "" {
		q = q.Where("goal = ?", f.Goal)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	return out, nil
}

func (r *Repository) FindTemplate(ctx context.Context, id uuid.UUID) (*models.PostTemplate, error) {
	var row models.PostTemplate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find template: %w", err)
	}
	return &row, nil
}

func (r *Repository) FindTemplateBySlug(ctx context.Context, slug string) (*models.PostTemplate, error) {
	var row models.PostTemplate
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find template by slug: %w", err)
	}
	return &row, nil
}

func (r *Repository) SaveTemplate(ctx context.Context, row *models.PostTemplate) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save template: %w", err)
	}
	return nil
}

func (r *Repository) CreateTemplate(ctx context.Context, row *models.PostTemplate) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	return r.SaveTemplate(ctx, row)
}

func (r *Repository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.PostTemplate{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete template: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) ListPosts(ctx context.Context, f PostFilter) ([]models.SocialPost, error) {
	var out []models.SocialPost
	q := r.db.WithContext(ctx).Order("COALESCE(scheduled_at, created_at) DESC")
	if f.Platform != "" {
		q = q.Where("platform = ?", f.Platform)
	}
	if f.Goal != "" {
		q = q.Where("goal = ?", f.Goal)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.ProjectID != nil {
		q = q.Where("project_id = ?", *f.ProjectID)
	}
	if f.CampaignID != nil {
		q = q.Where("campaign_id = ?", *f.CampaignID)
	}
	limit := f.Limit
	if limit <= 0 {
		limit = 100
	}
	q = q.Limit(limit)
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	return out, nil
}

func (r *Repository) FindPost(ctx context.Context, id uuid.UUID) (*models.SocialPost, error) {
	var row models.SocialPost
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find post: %w", err)
	}
	return &row, nil
}

func (r *Repository) SavePost(ctx context.Context, row *models.SocialPost) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save post: %w", err)
	}
	return nil
}

func (r *Repository) CreatePost(ctx context.Context, row *models.SocialPost) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	return r.SavePost(ctx, row)
}

func (r *Repository) DeletePost(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.SocialPost{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete post: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func ParseDate(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return &t
	}
	return nil
}
