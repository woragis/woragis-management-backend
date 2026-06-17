package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListVideos(ctx context.Context) ([]models.LeetcodeVideo, error) {
	var rows []models.LeetcodeVideo
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetVideo(ctx context.Context, id uuid.UUID) (*models.LeetcodeVideo, error) {
	var row models.LeetcodeVideo
	err := r.db.WithContext(ctx).Preload("Thumbnails", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("created_at DESC")
	}).First(&row, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound(apperrors.CodeContentVideoNotFound, apperrors.MsgContentVideoNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "load leetcode video", err)
	}
	return &row, nil
}

func (r *Repository) CreateVideo(ctx context.Context, row *models.LeetcodeVideo) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repository) UpdateVideo(ctx context.Context, row *models.LeetcodeVideo) error {
	result := r.db.WithContext(ctx).Save(row)
	if result.Error != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "update leetcode video", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound(apperrors.CodeContentVideoNotFound, apperrors.MsgContentVideoNotFound)
	}
	return nil
}

func (r *Repository) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("video_id = ?", id).Delete(&models.ContentThumbnail{}).Error; err != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "delete thumbnails", err)
	}
	result := r.db.WithContext(ctx).Delete(&models.LeetcodeVideo{}, "id = ?", id)
	if result.Error != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "delete leetcode video", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound(apperrors.CodeContentVideoNotFound, apperrors.MsgContentVideoNotFound)
	}
	return nil
}

func (r *Repository) ListThumbnails(ctx context.Context, videoID uuid.UUID) ([]models.ContentThumbnail, error) {
	var rows []models.ContentThumbnail
	err := r.db.WithContext(ctx).Where("video_id = ?", videoID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetThumbnail(ctx context.Context, videoID, id uuid.UUID) (*models.ContentThumbnail, error) {
	var row models.ContentThumbnail
	err := r.db.WithContext(ctx).First(&row, "id = ? AND video_id = ?", id, videoID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound(apperrors.CodeContentThumbnailNotFound, apperrors.MsgContentThumbnailNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "load thumbnail", err)
	}
	return &row, nil
}

func (r *Repository) CreateThumbnail(ctx context.Context, row *models.ContentThumbnail) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repository) UpdateThumbnail(ctx context.Context, row *models.ContentThumbnail) error {
	result := r.db.WithContext(ctx).Save(row)
	if result.Error != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "update thumbnail", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound(apperrors.CodeContentThumbnailNotFound, apperrors.MsgContentThumbnailNotFound)
	}
	return nil
}

func (r *Repository) DeleteThumbnail(ctx context.Context, videoID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.ContentThumbnail{}, "id = ? AND video_id = ?", id, videoID)
	if result.Error != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "delete thumbnail", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound(apperrors.CodeContentThumbnailNotFound, apperrors.MsgContentThumbnailNotFound)
	}
	return nil
}

func (r *Repository) FindThumbnailByCreativesJob(ctx context.Context, jobID uuid.UUID) (*models.ContentThumbnail, error) {
	var row models.ContentThumbnail
	err := r.db.WithContext(ctx).First(&row, "creatives_job_id = ?", jobID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound(apperrors.CodeContentThumbnailNotFound, apperrors.MsgContentThumbnailNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "find thumbnail by creatives job", err)
	}
	return &row, nil
}

func (r *Repository) ListTemplates(ctx context.Context) ([]models.ContentPromptTemplate, error) {
	var rows []models.ContentPromptTemplate
	err := r.db.WithContext(ctx).Where("channel_slug = ?", "leetcode").Order("is_default DESC, name ASC").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetTemplate(ctx context.Context, id uuid.UUID) (*models.ContentPromptTemplate, error) {
	var row models.ContentPromptTemplate
	err := r.db.WithContext(ctx).First(&row, "id = ? AND channel_slug = ?", id, "leetcode").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound(apperrors.CodeContentTemplateNotFound, apperrors.MsgContentTemplateNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "load template", err)
	}
	return &row, nil
}

func (r *Repository) CreateTemplate(ctx context.Context, row *models.ContentPromptTemplate) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repository) UpdateTemplate(ctx context.Context, row *models.ContentPromptTemplate) error {
	result := r.db.WithContext(ctx).Save(row)
	if result.Error != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "update template", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound(apperrors.CodeContentTemplateNotFound, apperrors.MsgContentTemplateNotFound)
	}
	return nil
}

func (r *Repository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.ContentPromptTemplate{}, "id = ? AND channel_slug = ?", id, "leetcode")
	if result.Error != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "delete template", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound(apperrors.CodeContentTemplateNotFound, apperrors.MsgContentTemplateNotFound)
	}
	return nil
}

func (r *Repository) ClearDefaultTemplates(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&models.ContentPromptTemplate{}).
		Where("channel_slug = ?", "leetcode").
		Update("is_default", false).Error
}
