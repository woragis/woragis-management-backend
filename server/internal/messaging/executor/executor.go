package executor

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/agentworkerclient"
	"github.com/woragis/management/backend/server/internal/apperrors"
	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
	msgtemplaterender "github.com/woragis/management/backend/server/internal/messaging/templaterender"
	"github.com/woragis/management/backend/server/internal/models"
	"github.com/woragis/management/backend/server/internal/telegramworkerclient"
	"github.com/woragis/management/backend/server/internal/whatsappworkerclient"
)

type Executor struct {
	messaging *messagingsvc.Service
	content   *contentsvc.Service
	whatsapp  *whatsappworkerclient.Client
	telegram  *telegramworkerclient.Client
	agent     *agentworkerclient.Client
	renderer  *msgtemplaterender.Engine
}

func New(
	messaging *messagingsvc.Service,
	content *contentsvc.Service,
	whatsapp *whatsappworkerclient.Client,
	telegram *telegramworkerclient.Client,
	agent *agentworkerclient.Client,
	renderer *msgtemplaterender.Engine,
) *Executor {
	return &Executor{
		messaging: messaging,
		content:   content,
		whatsapp:  whatsapp,
		telegram:  telegram,
		agent:     agent,
		renderer:  renderer,
	}
}

type ExecuteResult struct {
	JobID      uuid.UUID `json:"jobId"`
	Skipped    bool      `json:"skipped"`
	SkipReason string    `json:"skipReason,omitempty"`
	Message    string    `json:"message,omitempty"`
	DeliveryID uuid.UUID `json:"deliveryId,omitempty"`
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

	message, templateSlug, externalRef, tpl, renderData, skip, skipReason, err := e.resolveMessage(ctx, job, dest)
	if err != nil {
		return nil, err
	}
	if skip {
		if err := e.finishSkipped(ctx, job, dest, templateSlug, skipReason, externalRef); err != nil {
			return nil, err
		}
		return &ExecuteResult{JobID: jobID, Skipped: true, SkipReason: skipReason}, nil
	}
	if strings.TrimSpace(message) == "" {
		if err := e.finishSkipped(ctx, job, dest, templateSlug, "empty message", externalRef); err != nil {
			return nil, err
		}
		return &ExecuteResult{JobID: jobID, Skipped: true, SkipReason: "empty message"}, nil
	}

	if tpl != nil && tpl.ComposeMode == models.ComposeModeAIAssisted {
		message, err = e.composeWithAgent(ctx, tpl, message, renderData, dest)
		if err != nil {
			_, _ = e.recordDelivery(ctx, job, dest, templateSlug, message, models.DeliveryStatusFailed, err.Error(), externalRef)
			return nil, apperrors.InternalErr(apperrors.CodeInternal, err.Error())
		}
	}

	if err := e.send(ctx, dest, message, job.ProgramAction, externalRef, templateSlug); err != nil {
		_, _ = e.recordDelivery(ctx, job, dest, templateSlug, message, models.DeliveryStatusFailed, err.Error(), externalRef)
		return nil, apperrors.InternalErr(apperrors.CodeInternal, err.Error())
	}

	delivery, err := e.recordDelivery(ctx, job, dest, templateSlug, message, models.DeliveryStatusSent, "", externalRef)
	if err != nil {
		return nil, err
	}
	if err := e.messaging.MarkJobRun(ctx, job, time.Now().UTC()); err != nil {
		return nil, err
	}

	e.patchLeetcodeSent(ctx, job, externalRef)

	return &ExecuteResult{
		JobID:      jobID,
		Message:    message,
		DeliveryID: delivery.ID,
	}, nil
}

func (e *Executor) composeWithAgent(ctx context.Context, tpl *models.MessageTemplate, staticBody string, data map[string]string, dest *models.ChannelDestination) (string, error) {
	if e.agent == nil || !e.agent.Enabled() {
		return staticBody, nil
	}
	dataAny := map[string]any{}
	for k, v := range data {
		dataAny[k] = v
	}
	return e.agent.Compose(ctx, agentworkerclient.ComposeRequest{
		TemplateBody: tpl.Body,
		ComposeMode:  tpl.ComposeMode,
		Data:         dataAny,
		DestinationContext: map[string]any{
			"channel":       dest.Channel,
			"destinationId": dest.ID.String(),
			"name":          dest.Name,
			"aiPromptHint":  tpl.AIPromptHint,
			"staticPreview": staticBody,
		},
	})
}

func (e *Executor) patchLeetcodeSent(ctx context.Context, job *models.ScheduledJob, externalRef string) {
	if externalRef == "" || e.content == nil {
		return
	}
	if !isLeetcodeDispatchType(strings.TrimSpace(job.ProgramAction)) {
		return
	}
	vid, err := uuid.Parse(externalRef)
	if err != nil {
		return
	}
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

func (e *Executor) resolveMessage(ctx context.Context, job *models.ScheduledJob, dest *models.ChannelDestination) (message, templateSlug, externalRef string, tpl *models.MessageTemplate, renderData map[string]string, skip bool, skipReason string, err error) {
	action := normalizeProgramAction(job)
	if isLeetcodeDispatchType(action) && strings.TrimSpace(job.TemplateSlug) == "" {
		msg, slug, ref, t, sk, reason, e := e.resolveLeetcodeLegacy(ctx, action)
		return msg, slug, ref, t, nil, sk, reason, e
	}
	if strings.TrimSpace(job.TemplateSlug) == "" {
		return "", "", "", nil, nil, true, "no template or program action", nil
	}
	tpl, err = e.messaging.FindTemplateBySlug(ctx, job.TemplateSlug, dest.ID)
	if err != nil {
		return "", job.TemplateSlug, "", nil, nil, true, "template not found", nil
	}
	templateSlug = tpl.Slug

	if e.renderer != nil {
		res, err := e.renderer.Render(ctx, msgtemplaterender.RenderInput{Template: tpl, Job: job})
		if err != nil {
			return "", templateSlug, "", tpl, nil, false, "", err
		}
		if res.Skipped {
			return "", templateSlug, res.ExternalRef, tpl, nil, true, res.SkipReason, nil
		}
		return res.Body, templateSlug, res.ExternalRef, tpl, res.Data, false, "", nil
	}
	return tpl.Body, templateSlug, "", tpl, nil, false, "", nil
}

func (e *Executor) resolveLeetcodeLegacy(ctx context.Context, dispatchType string) (message, templateSlug, externalRef string, tpl *models.MessageTemplate, skip bool, skipReason string, err error) {
	if e.content == nil {
		return "", "", "", nil, true, "content service unavailable", nil
	}
	d, err := e.content.Dispatch(ctx, dispatchType, "")
	if err != nil {
		return "", "", "", nil, false, "", err
	}
	if d.Skip {
		return "", d.TemplateSlug, d.VideoID, nil, true, d.SkipReason, nil
	}
	return d.Message, d.TemplateSlug, d.VideoID, nil, false, "", nil
}

func normalizeProgramAction(job *models.ScheduledJob) string {
	action := strings.TrimSpace(job.ProgramAction)
	if action == "" && strings.HasPrefix(strings.TrimSpace(job.TemplateSlug), "leetcode/") {
		action = strings.TrimPrefix(strings.TrimSpace(job.TemplateSlug), "leetcode/")
	}
	if strings.HasPrefix(action, "leetcode/") {
		action = strings.TrimPrefix(action, "leetcode/")
	}
	return action
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
			Message:       message,
			ExternalID:    dest.ExternalID,
			DestinationID: dest.ID.String(),
			Type:          dispatchType,
			VideoID:       videoID,
			TemplateSlug:  templateSlug,
		})
	case models.ChannelTelegram:
		if e.telegram == nil || !e.telegram.Enabled() {
			return apperrors.Unavailable(apperrors.CodeInternal, "Telegram worker not configured.")
		}
		return e.telegram.Send(ctx, telegramworkerclient.SendRequest{
			Message:    message,
			ExternalID: dest.ExternalID,
		})
	default:
		return apperrors.Invalid(apperrors.CodeInternal, "Unsupported channel.")
	}
}

// PreviewTemplate renders a template without sending (admin preview).
func (e *Executor) PreviewTemplate(ctx context.Context, tpl *models.MessageTemplate, job *models.ScheduledJob) (*msgtemplaterender.RenderResult, error) {
	if e.renderer == nil {
		return &msgtemplaterender.RenderResult{Body: tpl.Body}, nil
	}
	return e.renderer.Render(ctx, msgtemplaterender.RenderInput{Template: tpl, Job: job})
}
