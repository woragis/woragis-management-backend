package reminder

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
	"github.com/woragis/management/backend/server/internal/models"
	presencesvc "github.com/woragis/management/backend/server/internal/presence/service"
	"github.com/woragis/management/backend/server/internal/whatsappworkerclient"
)

type Executor struct {
	presence    *presencesvc.Service
	messaging   *messagingsvc.Service
	devProjects *devprojectsvc.Service
	whatsapp    *whatsappworkerclient.Client
}

func New(
	presence *presencesvc.Service,
	messaging *messagingsvc.Service,
	devProjects *devprojectsvc.Service,
	whatsapp *whatsappworkerclient.Client,
) *Executor {
	return &Executor{
		presence:    presence,
		messaging:   messaging,
		devProjects: devProjects,
		whatsapp:    whatsapp,
	}
}

type SendResult struct {
	PostID     uuid.UUID `json:"postId"`
	Skipped    bool      `json:"skipped"`
	SkipReason string    `json:"skipReason,omitempty"`
	Message    string    `json:"message,omitempty"`
}

func (e *Executor) SendReminder(ctx context.Context, postID uuid.UUID) (*SendResult, error) {
	settings, err := e.presence.GetSettings(ctx)
	if err != nil {
		return nil, err
	}
	if !settings.RemindersEnabled {
		return &SendResult{PostID: postID, Skipped: true, SkipReason: "reminders disabled"}, nil
	}

	post, err := e.presence.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.Status != models.SocialPostStatusScheduled {
		return &SendResult{PostID: postID, Skipped: true, SkipReason: "post not scheduled"}, nil
	}
	if post.ReminderSentAt != nil {
		return &SendResult{PostID: postID, Skipped: true, SkipReason: "reminder already sent"}, nil
	}
	if post.ScheduledAt == nil || post.ScheduledAt.After(time.Now().UTC()) {
		return &SendResult{PostID: postID, Skipped: true, SkipReason: "not due yet"}, nil
	}

	destID := settings.DefaultDestinationID
	if post.NotifyDestinationID != nil {
		destID = post.NotifyDestinationID
	}
	if destID == nil {
		return &SendResult{PostID: postID, Skipped: true, SkipReason: "no notify destination configured"}, nil
	}

	dest, err := e.messaging.GetDestination(ctx, *destID)
	if err != nil {
		return nil, err
	}
	if !dest.Active {
		return &SendResult{PostID: postID, Skipped: true, SkipReason: "destination inactive"}, nil
	}
	if dest.Channel != models.ChannelWhatsApp {
		return &SendResult{PostID: postID, Skipped: true, SkipReason: "only whatsapp reminders supported"}, nil
	}

	message := formatReminderMessage(post, e.projectName(ctx, post.ProjectID))
	if strings.TrimSpace(message) == "" {
		return &SendResult{PostID: postID, Skipped: true, SkipReason: "empty message"}, nil
	}

	if e.whatsapp == nil || !e.whatsapp.Enabled() {
		return nil, apperrors.Unavailable(apperrors.CodeWhatsappWorkerUnavailable, apperrors.MsgWhatsappWorkerUnavailable)
	}
	if err := e.whatsapp.Send(ctx, whatsappworkerclient.SendRequest{
		Message:       message,
		ExternalID:    dest.ExternalID,
		DestinationID: dest.ID.String(),
		Type:          "presence_reminder",
		TemplateSlug:  "presence/reminder",
	}); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := e.presence.MarkReminderSent(ctx, postID, now); err != nil {
		return nil, err
	}

	delivery := &models.MessageDelivery{
		DestinationID: dest.ID,
		Channel:       dest.Channel,
		ExternalID:    dest.ExternalID,
		TemplateSlug:  "presence/reminder",
		Body:          message,
		Status:        "sent",
		ExternalRef:   postID.String(),
		SentAt:        now,
	}
	_ = e.messaging.RecordDelivery(ctx, delivery)

	return &SendResult{PostID: postID, Message: message}, nil
}

func (e *Executor) projectName(ctx context.Context, projectID *uuid.UUID) string {
	if projectID == nil || e.devProjects == nil {
		return ""
	}
	p, err := e.devProjects.GetByID(ctx, *projectID)
	if err != nil {
		return ""
	}
	return p.Name
}

func formatReminderMessage(post *models.SocialPost, projectName string) string {
	platform := post.Platform
	switch post.Platform {
	case models.SocialPlatformLinkedIn:
		platform = "LinkedIn"
	case models.SocialPlatformTwitter:
		platform = "Twitter / X"
	case models.SocialPlatformReddit:
		platform = "Reddit"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "📣 Post agendado — %s\n", platform)
	if projectName != "" {
		fmt.Fprintf(&b, "Projeto: %s\n", projectName)
	}
	if post.Title != "" {
		fmt.Fprintf(&b, "Título: %s\n", post.Title)
	}
	b.WriteString("---\n")
	parts := []string{}
	if post.Hook != "" {
		parts = append(parts, post.Hook)
	}
	if post.Body != "" {
		parts = append(parts, post.Body)
	}
	if post.CTA != "" {
		parts = append(parts, post.CTA)
	}
	b.WriteString(strings.Join(parts, "\n\n"))
	b.WriteString("\n---\nPublique manualmente e marque como publicado no admin.")
	return b.String()
}
