package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

func (r *Repository) GetSettings(ctx context.Context) (*models.LeetcodeChannelSettings, error) {
	var row models.LeetcodeChannelSettings
	err := r.db.WithContext(ctx).First(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound(apperrors.CodeContentSettingsNotFound, apperrors.MsgContentSettingsNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "load leetcode settings", err)
	}
	return &row, nil
}

func (r *Repository) EnsureSettings(ctx context.Context) (*models.LeetcodeChannelSettings, error) {
	var row models.LeetcodeChannelSettings
	err := r.db.WithContext(ctx).First(&row).Error
	if err == nil {
		return &row, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	row = models.LeetcodeChannelSettings{ID: uuid.New()}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repository) UpdateSettings(ctx context.Context, row *models.LeetcodeChannelSettings) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repository) ListWhatsappTemplates(ctx context.Context) ([]models.WhatsappMessageTemplate, error) {
	var rows []models.WhatsappMessageTemplate
	err := r.db.WithContext(ctx).Where("channel_slug = ?", "leetcode").Order("slug ASC").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetWhatsappTemplateBySlug(ctx context.Context, slug string) (*models.WhatsappMessageTemplate, error) {
	var row models.WhatsappMessageTemplate
	err := r.db.WithContext(ctx).First(&row, "channel_slug = ? AND slug = ?", "leetcode", slug).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound(apperrors.CodeContentWhatsappTplNotFound, apperrors.MsgContentWhatsappTplNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "load whatsapp template", err)
	}
	return &row, nil
}

func (r *Repository) GetWhatsappTemplate(ctx context.Context, id uuid.UUID) (*models.WhatsappMessageTemplate, error) {
	var row models.WhatsappMessageTemplate
	err := r.db.WithContext(ctx).First(&row, "id = ? AND channel_slug = ?", id, "leetcode").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound(apperrors.CodeContentWhatsappTplNotFound, apperrors.MsgContentWhatsappTplNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "load whatsapp template", err)
	}
	return &row, nil
}

func (r *Repository) CreateWhatsappTemplate(ctx context.Context, row *models.WhatsappMessageTemplate) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repository) UpdateWhatsappTemplate(ctx context.Context, row *models.WhatsappMessageTemplate) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repository) DeleteWhatsappTemplate(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.WhatsappMessageTemplate{}, "id = ? AND channel_slug = ?", id, "leetcode")
	if result.RowsAffected == 0 {
		return apperrors.NotFound(apperrors.CodeContentWhatsappTplNotFound, apperrors.MsgContentWhatsappTplNotFound)
	}
	return result.Error
}

func (r *Repository) WhatsappTemplateExists(ctx context.Context, slug string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.WhatsappMessageTemplate{}).
		Where("channel_slug = ? AND slug = ?", "leetcode", slug).
		Count(&n).Error
	return n > 0, err
}

func (r *Repository) GetVideoByProblemDate(ctx context.Context, day time.Time) (*models.LeetcodeVideo, error) {
	var row models.LeetcodeVideo
	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Where("problem_date >= ? AND problem_date < ? AND whatsapp_enabled = ?", dayStart, dayEnd, true).
		Order("created_at DESC").
		First(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound(apperrors.CodeContentVideoNotFound, apperrors.MsgContentVideoNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "load video by date", err)
	}
	return &row, nil
}

func (r *Repository) ListVideosWithProblemSentBetween(ctx context.Context, from, to time.Time) ([]models.LeetcodeVideo, error) {
	var rows []models.LeetcodeVideo
	err := r.db.WithContext(ctx).
		Where("whatsapp_problem_sent_at IS NOT NULL AND whatsapp_problem_sent_at >= ? AND whatsapp_problem_sent_at < ?", from, to).
		Order("series_number ASC, created_at ASC").
		Find(&rows).Error
	return rows, err
}
