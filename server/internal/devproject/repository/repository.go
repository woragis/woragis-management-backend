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

func (r *Repository) ListProjects(ctx context.Context) ([]models.Project, error) {
	var out []models.Project
	if err := r.db.WithContext(ctx).Order("display_order ASC, created_at DESC").Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return out, nil
}

func (r *Repository) FindProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var p models.Project
	err := r.db.WithContext(ctx).
		Preload("Links").
		Preload("Domains").
		Preload("Gallery", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Preload("Envs", func(db *gorm.DB) *gorm.DB {
			return db.Order("key ASC")
		}).
		Where("id = ?", id).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find project: %w", err)
	}
	return &p, nil
}

func (r *Repository) FindProjectByPublicSlug(ctx context.Context, slug string) (*models.Project, error) {
	var p models.Project
	err := r.db.WithContext(ctx).
		Preload("Links", "is_public = ?", true).
		Preload("Gallery", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Where("public_slug = ? AND is_public = ?", slug, true).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find project by slug: %w", err)
	}
	return &p, nil
}

func (r *Repository) ListPublicProjects(ctx context.Context, featuredOnly bool) ([]models.Project, error) {
	var out []models.Project
	q := r.db.WithContext(ctx).
		Preload("Links", "is_public = ?", true).
		Preload("Gallery", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Where("is_public = ?", true)
	if featuredOnly {
		q = q.Where("featured = ?", true)
	}
	if err := q.Order("display_order ASC, created_at DESC").Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list public projects: %w", err)
	}
	return out, nil
}

func (r *Repository) CreateProject(ctx context.Context, p *models.Project) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (r *Repository) SaveProject(ctx context.Context, p *models.Project) error {
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return fmt.Errorf("save project: %w", err)
	}
	return nil
}

func (r *Repository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectLink{}).Error; err != nil {
			return fmt.Errorf("delete links: %w", err)
		}
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectDomain{}).Error; err != nil {
			return fmt.Errorf("delete domains: %w", err)
		}
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectSecret{}).Error; err != nil {
			return fmt.Errorf("delete secrets: %w", err)
		}
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectGallery{}).Error; err != nil {
			return fmt.Errorf("delete gallery: %w", err)
		}
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectEnv{}).Error; err != nil {
			return fmt.Errorf("delete envs: %w", err)
		}
		res := tx.Where("id = ?", id).Delete(&models.Project{})
		if res.Error != nil {
			return fmt.Errorf("delete project: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *Repository) CreateLink(ctx context.Context, link *models.ProjectLink) error {
	if link.ID == uuid.Nil {
		link.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		return fmt.Errorf("create link: %w", err)
	}
	return nil
}

func (r *Repository) DeleteLink(ctx context.Context, projectID, linkID uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND project_id = ?", linkID, projectID).Delete(&models.ProjectLink{})
	if res.Error != nil {
		return fmt.Errorf("delete link: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) CreateDomain(ctx context.Context, d *models.ProjectDomain) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(d).Error; err != nil {
		return fmt.Errorf("create domain: %w", err)
	}
	return nil
}

func (r *Repository) DeleteDomain(ctx context.Context, projectID, domainID uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND project_id = ?", domainID, projectID).Delete(&models.ProjectDomain{})
	if res.Error != nil {
		return fmt.Errorf("delete domain: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) ListSecrets(ctx context.Context, projectID uuid.UUID) ([]models.ProjectSecret, error) {
	var out []models.ProjectSecret
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("name ASC").Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list secrets: %w", err)
	}
	return out, nil
}

func (r *Repository) FindSecret(ctx context.Context, projectID, secretID uuid.UUID) (*models.ProjectSecret, error) {
	var s models.ProjectSecret
	err := r.db.WithContext(ctx).Where("id = ? AND project_id = ?", secretID, projectID).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find secret: %w", err)
	}
	return &s, nil
}

func (r *Repository) CreateSecret(ctx context.Context, s *models.ProjectSecret) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(s).Error; err != nil {
		return fmt.Errorf("create secret: %w", err)
	}
	return nil
}

func (r *Repository) DeleteSecret(ctx context.Context, projectID, secretID uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND project_id = ?", secretID, projectID).Delete(&models.ProjectSecret{})
	if res.Error != nil {
		return fmt.Errorf("delete secret: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) CreateGalleryItem(ctx context.Context, item *models.ProjectGallery) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create gallery item: %w", err)
	}
	return nil
}

func (r *Repository) DeleteGalleryItem(ctx context.Context, projectID, itemID uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND project_id = ?", itemID, projectID).Delete(&models.ProjectGallery{})
	if res.Error != nil {
		return fmt.Errorf("delete gallery item: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
