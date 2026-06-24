package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
)

type UpdateSettingsInput struct {
	RemindersEnabled     *bool
	DefaultDestinationID *uuid.UUID
	DestinationSet       bool
}

func (s *Service) GetSettings(ctx context.Context) (*models.PresenceSettings, error) {
	row, err := s.repo.EnsureSettings(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load presence settings.", err)
	}
	return row, nil
}

func (s *Service) UpdateSettings(ctx context.Context, in UpdateSettingsInput) (*models.PresenceSettings, error) {
	row, err := s.repo.EnsureSettings(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load presence settings.", err)
	}
	if in.RemindersEnabled != nil {
		row.RemindersEnabled = *in.RemindersEnabled
	}
	if in.DestinationSet {
		row.DefaultDestinationID = in.DefaultDestinationID
	}
	if err := s.repo.SaveSettings(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update presence settings.", err)
	}
	return row, nil
}

func (s *Service) ListDueReminders(ctx context.Context, before time.Time, limit int) ([]models.SocialPost, error) {
	rows, err := s.repo.ListDueReminders(ctx, before, limit)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load due reminders.", err)
	}
	return rows, nil
}

func (s *Service) MarkReminderSent(ctx context.Context, postID uuid.UUID, at time.Time) error {
	row, err := s.GetPost(ctx, postID)
	if err != nil {
		return err
	}
	row.ReminderSentAt = &at
	if err := s.repo.SavePost(ctx, row); err != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to mark reminder sent.", err)
	}
	return nil
}
