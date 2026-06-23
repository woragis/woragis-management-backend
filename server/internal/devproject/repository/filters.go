package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

type ListFilter struct {
	Status         string
	Intent         string
	Monetization   string
	Maturity       string
	VisibilityGoal string
	Distribution   string
	IsPublic       *bool
	Featured       *bool
	Query          string
}

func (r *Repository) ListProjectsFiltered(ctx context.Context, f ListFilter) ([]models.Project, error) {
	var out []models.Project
	q := r.db.WithContext(ctx)
	if s := strings.TrimSpace(f.Status); s != "" {
		q = q.Where("status = ?", s)
	}
	if s := strings.TrimSpace(f.Intent); s != "" {
		q = q.Where("intent = ?", s)
	}
	if s := strings.TrimSpace(f.Monetization); s != "" {
		q = q.Where("monetization = ?", s)
	}
	if s := strings.TrimSpace(f.Maturity); s != "" {
		q = q.Where("maturity = ?", s)
	}
	if s := strings.TrimSpace(f.VisibilityGoal); s != "" {
		q = q.Where("visibility_goal = ?", s)
	}
	if s := strings.TrimSpace(f.Distribution); s != "" {
		q = q.Where("distribution::text ILIKE ?", "%\""+s+"\"%")
	}
	if f.IsPublic != nil {
		q = q.Where("is_public = ?", *f.IsPublic)
	}
	if f.Featured != nil {
		q = q.Where("featured = ?", *f.Featured)
	}
	if qstr := strings.TrimSpace(f.Query); qstr != "" {
		like := "%" + strings.ToLower(qstr) + "%"
		q = q.Where(
			"LOWER(name) LIKE ? OR LOWER(slug) LIKE ? OR LOWER(short_description) LIKE ? OR stack::text ILIKE ?",
			like, like, like, like,
		)
	}
	if err := q.Order("display_order ASC, created_at DESC").Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list projects filtered: %w", err)
	}
	return out, nil
}

func (r *Repository) CountProjects(ctx context.Context) (int64, error) {
	var n int64
	if err := r.db.WithContext(ctx).Model(&models.Project{}).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count projects: %w", err)
	}
	return n, nil
}

func (r *Repository) CountProjectsPublic(ctx context.Context) (int64, error) {
	var n int64
	if err := r.db.WithContext(ctx).Model(&models.Project{}).Where("is_public = ?", true).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count public projects: %w", err)
	}
	return n, nil
}

func (r *Repository) CountProjectsActive(ctx context.Context) (int64, error) {
	var n int64
	if err := r.db.WithContext(ctx).Model(&models.Project{}).Where("status = ?", "active").Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count active projects: %w", err)
	}
	return n, nil
}

func (r *Repository) ListSecretsExpiringBefore(ctx context.Context, before time.Time) ([]models.ProjectSecret, error) {
	var out []models.ProjectSecret
	err := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at <= ?", before).
		Order("expires_at ASC").
		Limit(20).
		Find(&out).Error
	if err != nil {
		return nil, fmt.Errorf("list expiring secrets: %w", err)
	}
	return out, nil
}

func (r *Repository) ListDomainsExpiringBefore(ctx context.Context, before time.Time) ([]models.ProjectDomain, error) {
	var out []models.ProjectDomain
	err := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at <= ?", before).
		Order("expires_at ASC").
		Limit(20).
		Find(&out).Error
	if err != nil {
		return nil, fmt.Errorf("list expiring domains: %w", err)
	}
	return out, nil
}

func (r *Repository) ListEnvs(ctx context.Context, projectID uuid.UUID) ([]models.ProjectEnv, error) {
	var out []models.ProjectEnv
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("key ASC").Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list envs: %w", err)
	}
	return out, nil
}

func (r *Repository) CreateEnv(ctx context.Context, e *models.ProjectEnv) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return fmt.Errorf("create env: %w", err)
	}
	return nil
}

func (r *Repository) DeleteEnv(ctx context.Context, projectID, envID uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND project_id = ?", envID, projectID).Delete(&models.ProjectEnv{})
	if res.Error != nil {
		return fmt.Errorf("delete env: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
