package executor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/models"
)

func (e *Executor) recordDelivery(ctx context.Context, job *models.ScheduledJob, dest *models.ChannelDestination, templateSlug, body, status, errMsg, externalRef string) (*models.MessageDelivery, error) {
	var jobID *uuid.UUID
	if job != nil {
		jobID = &job.ID
	}
	delivery := &models.MessageDelivery{
		JobID:         jobID,
		DestinationID: dest.ID,
		Channel:       dest.Channel,
		ExternalID:    dest.ExternalID,
		TemplateSlug:  templateSlug,
		Body:          body,
		Status:        status,
		ErrorMessage:  errMsg,
		ExternalRef:   externalRef,
		SentAt:        time.Now().UTC(),
	}
	if err := e.messaging.RecordDelivery(ctx, delivery); err != nil {
		return nil, err
	}
	return delivery, nil
}

func (e *Executor) finishSkipped(ctx context.Context, job *models.ScheduledJob, dest *models.ChannelDestination, templateSlug, skipReason, externalRef string) error {
	_, err := e.recordDelivery(ctx, job, dest, templateSlug, skipReason, models.DeliveryStatusSkipped, "", externalRef)
	if err != nil {
		return err
	}
	return e.messaging.MarkJobRun(ctx, job, time.Now().UTC())
}
