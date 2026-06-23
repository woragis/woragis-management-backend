package executor

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
	"github.com/woragis/management/backend/server/internal/models"
	"github.com/woragis/management/backend/server/internal/whatsappworkerclient"
)

type Executor struct {
	messaging *messagingsvc.Service
	content   *contentsvc.Service
	whatsapp  *whatsappworkerclient.Client
}

func New(messaging *messagingsvc.Service, content *contentsvc.Service, whatsapp *whatsappworkerclient.Client) *Executor {
	return &Executor{messaging: messaging, content: content, whatsapp: whatsapp}
}

type ExecuteResult struct {
	JobID       uuid.UUID `json:"jobId"`
	Skipped     bool      `json:"skipped"`
	SkipReason  string    `json:"skipReason,omitempty"`
	Message     string    `json:"message,omitempty"`
	DeliveryID  uuid.UUID `json:"deliveryId,omitempty"`
}

func (e *Executor) ExecuteJob(ctx context.Context, jobID uuid.UUID) (*ExecuteResult, error) {
	job, err := e.messaging.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if !job.Enabled {
		return &ExecuteResult{JobID: jobID, Skipped: true, SkipReason: "job disabled"}, nil
	}

	dest, err := e.messaging.GetDestination(ctx, job.DestinationID)
	if err != nil {
		return nil, err
	}
	if !dest.Active {
		return &ExecuteResult{JobID: jobID, Skipped: true, SkipReason: "destination inactive"}, nil
	}

	message, templateSlug, externalRef, skip, skipReason, err := e.resolveMessage(ctx, job, dest)
	if err != nil {
		return nil, err
	}
	if skip {
		_ = e.messaging.MarkJobRun(ctx, job, time.Now().UTC())
		return &ExecuteResult{JobID: jobID, Skipped: true, SkipReason: skipReason}, nil
	}
	if strings.TrimSpace(message) == "" {
		return &ExecuteResult{JobID: jobID, Skipped: true, SkipReason: "empty message"}, nil
	}

	status := "sent"
	var errMsg string
	if err := e.send(ctx, dest, message, job.ProgramAction, externalRef, templateSlug); err != nil {
		status = "failed"
		errMsg = err.Error()
	}

	delivery := &models.MessageDelivery{
		JobID:         &job.ID,
		DestinationID: dest.ID,
		Channel:       dest.Channel,
		ExternalID:    dest.ExternalID,
		TemplateSlug:  templateSlug,
		Body:          message,
		Status:        status,
		ErrorMessage:  errMsg,
		ExternalRef:   externalRef,
		SentAt:        time.Now().UTC(),
	}
	if err := e.messaging.RecordDelivery(ctx, delivery); err != nil {
		return nil, err
	}
	_ = e.messaging.MarkJobRun(ctx, job, time.Now().UTC())

	if status == "failed" {
		return nil, apperrors.InternalErr(apperrors.CodeInternal, errMsg)
	}
	if externalRef != "" && e.content != nil && isLeetcodeDispatchType(strings.TrimSpace(job.ProgramAction)) {
		if vid, err := uuid.Parse(externalRef); err == nil {
			patch := contentsvc.WhatsappStatusPatch{}
			switch strings.TrimSpace(job.ProgramAction) {
			case "problem":
				patch.ProblemSent = true
			case "discussion":
				patch.DiscussionSent = true
			case "solution":
				patch.SolutionSent = true
			}
			_ = e.content.PatchWhatsappStatus(ctx, vid, patch)
		}
	}
	return &ExecuteResult{
		JobID:      jobID,
		Message:    message,
		DeliveryID: delivery.ID,
	}, nil
}

func (e *Executor) resolveMessage(ctx context.Context, job *models.ScheduledJob, dest *models.ChannelDestination) (message, templateSlug, externalRef string, skip bool, skipReason string, err error) {
	action := strings.TrimSpace(job.ProgramAction)
	if action == "" && strings.HasPrefix(strings.TrimSpace(job.TemplateSlug), "leetcode/") {
		action = strings.TrimPrefix(strings.TrimSpace(job.TemplateSlug), "leetcode/")
	}
	if strings.HasPrefix(action, "leetcode/") {
		action = strings.TrimPrefix(action, "leetcode/")
	}

	switch {
	case strings.HasPrefix(action, "leetcode") || isLeetcodeDispatchType(action):
		dispatchType := action
		if strings.Contains(action, "/") {
			parts := strings.SplitN(action, "/", 2)
			if len(parts) == 2 {
				dispatchType = parts[1]
			}
		}
		if e.content == nil {
			return "", "", "", true, "content service unavailable", nil
		}
		d, err := e.content.Dispatch(ctx, dispatchType, "")
		if err != nil {
			return "", "", "", false, "", err
		}
		if d.Skip {
			return "", d.TemplateSlug, d.VideoID, true, d.SkipReason, nil
		}
		return d.Message, d.TemplateSlug, d.VideoID, false, "", nil
	default:
		if job.TemplateSlug == "" {
			return "", "", "", true, "no template or program action", nil
		}
		tpl, err := e.messaging.ListTemplates(ctx, messagingsvc.TemplateFilter{
			ProgramSlug:   "custom",
			DestinationID: &dest.ID,
			ActiveOnly:    true,
		})
		if err != nil {
			return "", "", "", false, "", err
		}
		for _, t := range tpl {
			if t.Slug == job.TemplateSlug {
				return t.Body, t.Slug, "", false, "", nil
			}
		}
		all, err := e.messaging.ListTemplates(ctx, messagingsvc.TemplateFilter{ActiveOnly: true})
		if err != nil {
			return "", "", "", false, "", err
		}
		for _, t := range all {
			if t.Slug == job.TemplateSlug {
				return t.Body, t.Slug, "", false, "", nil
			}
		}
		return "", job.TemplateSlug, "", true, "template not found", nil
	}
}

func isLeetcodeDispatchType(t string) bool {
	switch t {
	case "problem", "discussion", "solution", "weekly":
		return true
	default:
		return false
	}
}

func (e *Executor) send(ctx context.Context, dest *models.ChannelDestination, message, dispatchType, videoID, templateSlug string) error {
	switch dest.Channel {
	case models.ChannelWhatsApp:
		if e.whatsapp == nil || !e.whatsapp.Enabled() {
			return apperrors.Unavailable(apperrors.CodeWhatsappWorkerUnavailable, apperrors.MsgWhatsappWorkerUnavailable)
		}
		return e.whatsapp.Send(ctx, whatsappworkerclient.SendRequest{
			Message:      message,
			ExternalID:   dest.ExternalID,
			DestinationID: dest.ID.String(),
			Type:         dispatchType,
			VideoID:      videoID,
			TemplateSlug: templateSlug,
		})
	case models.ChannelTelegram:
		return apperrors.Unavailable(apperrors.CodeInternal, "Telegram delivery not wired yet.")
	default:
		return apperrors.Invalid(apperrors.CodeInternal, "Unsupported channel.")
	}
}
