package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/woragis/management/backend/server/internal/agent/personality/cache"
	"github.com/woragis/management/backend/server/internal/agent/personality/repository"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

var allowedVoices = map[string]bool{
	"alloy": true, "ash": true, "ballad": true, "coral": true, "echo": true,
	"sage": true, "shimmer": true, "verse": true, "marin": true, "cedar": true,
	"onyx": true, "nova": true, "fable": true,
}

const maxSystemPromptExtra = 2000

type Service struct {
	repo  *repository.Repository
	cache cache.Store
}

func New(repo *repository.Repository, c cache.Store) *Service {
	if c == nil {
		c = cache.Noop{}
	}
	return &Service{repo: repo, cache: c}
}

type UpdateInput struct {
	AssistantName     *string
	GreetingMorning   *string
	GreetingAfternoon *string
	GreetingEvening   *string
	GreetingEnabled   *bool
	SystemPromptExtra *string
	VoiceID           *string
	Language          *string
	Timezone          *string
}

func (s *Service) Get(ctx context.Context) (*models.AgentPersonality, error) {
	if raw, err := s.cache.Get(ctx, cache.PersonalityKey); err == nil && len(raw) > 0 {
		var row models.AgentPersonality
		if json.Unmarshal(raw, &row) == nil {
			return &row, nil
		}
	}
	row, err := s.repo.EnsureDefault(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load agent personality.", err)
	}
	_ = s.writeCache(ctx, row)
	return row, nil
}

func (s *Service) Update(ctx context.Context, in UpdateInput) (*models.AgentPersonality, error) {
	row, err := s.Get(ctx)
	if err != nil {
		return nil, err
	}
	if in.AssistantName != nil {
		name := strings.TrimSpace(*in.AssistantName)
		if name == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Assistant name cannot be empty.")
		}
		row.AssistantName = name
		now := time.Now().UTC()
		row.AssistantNameSetAt = &now
	}
	if in.GreetingMorning != nil {
		row.GreetingMorning = strings.TrimSpace(*in.GreetingMorning)
	}
	if in.GreetingAfternoon != nil {
		row.GreetingAfternoon = strings.TrimSpace(*in.GreetingAfternoon)
	}
	if in.GreetingEvening != nil {
		row.GreetingEvening = strings.TrimSpace(*in.GreetingEvening)
	}
	if in.GreetingEnabled != nil {
		row.GreetingEnabled = *in.GreetingEnabled
	}
	if in.SystemPromptExtra != nil {
		extra := strings.TrimSpace(*in.SystemPromptExtra)
		if len(extra) > maxSystemPromptExtra {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "System prompt extra is too long.")
		}
		row.SystemPromptExtra = extra
	}
	if in.VoiceID != nil {
		voice := strings.ToLower(strings.TrimSpace(*in.VoiceID))
		if !allowedVoices[voice] {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Voice id is not allowed.")
		}
		row.VoiceID = voice
	}
	if in.Language != nil {
		lang := strings.TrimSpace(*in.Language)
		if lang != "" {
			row.Language = lang
		}
	}
	if in.Timezone != nil {
		tz := strings.TrimSpace(*in.Timezone)
		if tz != "" {
			row.Timezone = tz
		}
	}
	if err := s.repo.Save(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update agent personality.", err)
	}
	_ = s.writeCache(ctx, row)
	return row, nil
}

func (s *Service) Reset(ctx context.Context) (*models.AgentPersonality, error) {
	def := models.DefaultAgentPersonality()
	existing, err := s.repo.Find(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to reset agent personality.", err)
	}
	if existing != nil {
		def.CreatedAt = existing.CreatedAt
	}
	if err := s.repo.Save(ctx, &def); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to reset agent personality.", err)
	}
	_ = s.writeCache(ctx, &def)
	return &def, nil
}

func (s *Service) writeCache(ctx context.Context, row *models.AgentPersonality) error {
	b, err := json.Marshal(row)
	if err != nil {
		return err
	}
	return s.cache.Set(ctx, cache.PersonalityKey, b)
}
