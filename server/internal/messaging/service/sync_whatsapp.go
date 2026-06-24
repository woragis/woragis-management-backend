package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	msgtemplaterender "github.com/woragis/management/backend/server/internal/messaging/templaterender"
	"github.com/woragis/management/backend/server/internal/whatsappworkerclient"
	"github.com/woragis/management/backend/server/internal/telegramworkerclient"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WhatsAppSyncResult struct {
	Created   int                   `json:"created"`
	Updated   int                   `json:"updated"`
	Unchanged int                   `json:"unchanged"`
	Groups    []WhatsAppSyncedGroup `json:"groups"`
}

type WhatsAppSyncedGroup struct {
	JID     string `json:"jid"`
	Name    string `json:"name"`
	ID      string `json:"id,omitempty"`
	Created bool   `json:"created"`
	Updated bool   `json:"updated"`
}

type TelegramSyncResult struct {
	Created   int                    `json:"created"`
	Updated   int                    `json:"updated"`
	Unchanged int                    `json:"unchanged"`
	Chats     []TelegramSyncedChat   `json:"chats"`
}

type TelegramSyncedChat struct {
	ChatID  string `json:"chatId"`
	Name    string `json:"name"`
	ID      string `json:"id,omitempty"`
	Created bool   `json:"created"`
	Updated bool   `json:"updated"`
}

func (s *Service) ResolveDestination(ctx context.Context, channel, externalID string) (*models.ChannelDestination, error) {
	channel = strings.TrimSpace(strings.ToLower(channel))
	externalID = strings.TrimSpace(externalID)
	if channel == "" || externalID == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Channel and external id are required.")
	}
	row, err := s.repo.FindDestinationByExternal(ctx, channel, externalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Destination not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to resolve destination.", err)
	}
	return row, nil
}

func (s *Service) SyncWhatsAppDestinations(ctx context.Context, wa *whatsappworkerclient.Client) (*WhatsAppSyncResult, error) {
	if wa == nil || !wa.Enabled() {
		return nil, apperrors.Unavailable(apperrors.CodeWhatsappWorkerUnavailable, apperrors.MsgWhatsappWorkerUnavailable)
	}
	resp, err := wa.ListGroups(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to list WhatsApp groups.", err)
	}
	out := &WhatsAppSyncResult{}
	for _, g := range resp.Groups {
		jid := strings.TrimSpace(g.JID)
		if jid == "" {
			continue
		}
		name := strings.TrimSpace(g.Name)
		if name == "" {
			name = jid
		}
		existing, err := s.repo.FindDestinationByExternal(ctx, models.ChannelWhatsApp, jid)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to lookup destination.", err)
		}
		entry := WhatsAppSyncedGroup{JID: jid, Name: name}
		if existing != nil {
			entry.ID = existing.ID.String()
			if existing.Name != name {
				existing.Name = name
				if err := s.repo.SaveDestination(ctx, existing); err != nil {
					return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update destination.", err)
				}
				entry.Updated = true
				out.Updated++
			} else {
				out.Unchanged++
			}
			out.Groups = append(out.Groups, entry)
			continue
		}
		row := &models.ChannelDestination{
			Channel:    models.ChannelWhatsApp,
			ExternalID: jid,
			Name:       name,
			Description: "Synced from WhatsApp worker",
			Active:     false,
			Tags:       datatypes.JSON([]byte(`["whatsapp-sync"]`)),
		}
		if err := s.repo.CreateDestination(ctx, row); err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create destination.", err)
		}
		entry.ID = row.ID.String()
		entry.Created = true
		out.Created++
		out.Groups = append(out.Groups, entry)
	}
	return out, nil
}

func (s *Service) SyncTelegramDestinations(ctx context.Context, tg *telegramworkerclient.Client) (*TelegramSyncResult, error) {
	if tg == nil || !tg.Enabled() {
		return nil, apperrors.Unavailable(apperrors.CodeInternal, "Telegram worker not configured.")
	}
	resp, err := tg.ListChats(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to list Telegram chats.", err)
	}
	out := &TelegramSyncResult{}
	for _, chat := range resp.Chats {
		chatID := strings.TrimSpace(chat.ChatID)
		if chatID == "" {
			continue
		}
		name := strings.TrimSpace(chat.Name)
		if name == "" {
			name = chatID
		}
		existing, err := s.repo.FindDestinationByExternal(ctx, models.ChannelTelegram, chatID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to lookup destination.", err)
		}
		entry := TelegramSyncedChat{ChatID: chatID, Name: name}
		if existing != nil {
			entry.ID = existing.ID.String()
			if existing.Name != name {
				existing.Name = name
				if err := s.repo.SaveDestination(ctx, existing); err != nil {
					return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update destination.", err)
				}
				entry.Updated = true
				out.Updated++
			} else {
				out.Unchanged++
			}
			out.Chats = append(out.Chats, entry)
			continue
		}
		row := &models.ChannelDestination{
			Channel:     models.ChannelTelegram,
			ExternalID:  chatID,
			Name:        name,
			Description: "Synced from Telegram worker",
			Active:      false,
			Tags:        datatypes.JSON([]byte(`["telegram-sync"]`)),
		}
		if err := s.repo.CreateDestination(ctx, row); err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create destination.", err)
		}
		entry.ID = row.ID.String()
		entry.Created = true
		out.Created++
		out.Chats = append(out.Chats, entry)
	}
	return out, nil
}

func failureBackoff(count int) time.Duration {
	if count < 1 {
		count = 1
	}
	minutes := 5
	for i := 1; i < count && minutes < 60; i++ {
		minutes *= 2
	}
	if minutes > 60 {
		minutes = 60
	}
	return time.Duration(minutes) * time.Minute
}

func (s *Service) MarkJobFailure(ctx context.Context, job *models.ScheduledJob, failedAt time.Time) error {
	job.FailureCount++
	job.NextRunAt = ptrTime(failedAt.Add(failureBackoff(job.FailureCount)))
	return s.repo.SaveJob(ctx, job)
}

func (s *Service) MarkJobRun(ctx context.Context, job *models.ScheduledJob, ranAt time.Time) error {
	job.FailureCount = 0
	job.LastRunAt = &ranAt
	next, err := computeNextRun(job.CronExpr, job.Timezone, ranAt)
	if err != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to compute next run.", err)
	}
	job.NextRunAt = &next
	return s.repo.SaveJob(ctx, job)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

// EnsureLeetcodeTemplates copies content WhatsappMessageTemplate rows into MessageTemplate when missing.
func (s *Service) EnsureLeetcodeTemplates(ctx context.Context, source []models.WhatsappMessageTemplate) error {
	bindings := msgtemplaterender.DefaultBindings("leetcode")
	bindingsJSON, _ := json.Marshal(bindings)
	for _, src := range source {
		if strings.TrimSpace(src.ChannelSlug) != "" && src.ChannelSlug != "leetcode" {
			continue
		}
		slug := strings.TrimSpace(src.Slug)
		if slug == "" {
			continue
		}
		if _, err := s.repo.FindTemplateBySlug(ctx, "leetcode", slug, nil); err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.InternalCause(apperrors.CodeInternal, "Failed to check template.", err)
		}
		row := &models.MessageTemplate{
			ProgramSlug: "leetcode",
			Slug:        slug,
			Name:        strings.TrimSpace(src.Name),
			Body:        src.Body,
			ComposeMode: models.ComposeModeStatic,
			Bindings:    datatypes.JSON(bindingsJSON),
			Active:      true,
		}
		if row.Name == "" {
			row.Name = slug
		}
		if err := s.repo.CreateTemplate(ctx, row); err != nil {
			return apperrors.InternalCause(apperrors.CodeInternal, "Failed to seed leetcode template.", err)
		}
	}
	return nil
}
