package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/messaging/repository"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

type DestinationFilter struct {
	Channel    string
	ActiveOnly bool
	Query      string
}

type CreateDestinationInput struct {
	Channel          string
	ExternalID       string
	Name             string
	Description      string
	Responsibilities string
	Tags             []string
	Metadata         map[string]any
	Active           bool
}

type UpdateDestinationInput struct {
	Channel          *string
	ExternalID       *string
	Name             *string
	Description      *string
	Responsibilities *string
	Tags             []string
	TagsSet          bool
	Metadata         map[string]any
	MetadataSet      bool
	Active           *bool
}

func (s *Service) ListDestinations(ctx context.Context, f DestinationFilter) ([]models.ChannelDestination, error) {
	rows, err := s.repo.ListDestinations(ctx, repository.DestinationFilter(f))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load destinations.", err)
	}
	return rows, nil
}

func (s *Service) GetDestination(ctx context.Context, id uuid.UUID) (*models.ChannelDestination, error) {
	row, err := s.repo.FindDestination(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Destination not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load destination.", err)
	}
	return row, nil
}

func (s *Service) CreateDestination(ctx context.Context, in CreateDestinationInput) (*models.ChannelDestination, error) {
	channel := normalizeChannel(in.Channel)
	if channel == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Channel must be whatsapp or telegram.")
	}
	externalID := strings.TrimSpace(in.ExternalID)
	name := strings.TrimSpace(in.Name)
	if externalID == "" || name == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "External id and name are required.")
	}
	row := &models.ChannelDestination{
		Channel:          channel,
		ExternalID:       externalID,
		Name:             name,
		Description:      strings.TrimSpace(in.Description),
		Responsibilities: strings.TrimSpace(in.Responsibilities),
		Tags:             jsonSlice(in.Tags),
		Metadata:         jsonMap(in.Metadata),
		Active:           in.Active,
	}
	if err := s.repo.CreateDestination(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create destination.", err)
	}
	return row, nil
}

func (s *Service) UpdateDestination(ctx context.Context, id uuid.UUID, in UpdateDestinationInput) (*models.ChannelDestination, error) {
	row, err := s.GetDestination(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Channel != nil {
		ch := normalizeChannel(*in.Channel)
		if ch == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Channel must be whatsapp or telegram.")
		}
		row.Channel = ch
	}
	if in.ExternalID != nil {
		ext := strings.TrimSpace(*in.ExternalID)
		if ext == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "External id cannot be empty.")
		}
		row.ExternalID = ext
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Name cannot be empty.")
		}
		row.Name = name
	}
	if in.Description != nil {
		row.Description = strings.TrimSpace(*in.Description)
	}
	if in.Responsibilities != nil {
		row.Responsibilities = strings.TrimSpace(*in.Responsibilities)
	}
	if in.TagsSet {
		row.Tags = jsonSlice(in.Tags)
	}
	if in.MetadataSet {
		row.Metadata = jsonMap(in.Metadata)
	}
	if in.Active != nil {
		row.Active = *in.Active
	}
	if err := s.repo.SaveDestination(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update destination.", err)
	}
	return row, nil
}

func (s *Service) DeleteDestination(ctx context.Context, id uuid.UUID) error {
	row, err := s.GetDestination(ctx, id)
	if err != nil {
		return err
	}
	row.Active = false
	if err := s.repo.SaveDestination(ctx, row); err != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to deactivate destination.", err)
	}
	return nil
}

type TemplateFilter struct {
	DestinationID *uuid.UUID
	ProgramSlug   string
	ActiveOnly    bool
}

type CreateTemplateInput struct {
	DestinationID *uuid.UUID
	ProgramSlug   string
	Slug          string
	Name          string
	Body          string
	ComposeMode   string
	AIPromptHint  string
	Bindings      map[string]string
	Active        bool
}

type UpdateTemplateInput struct {
	DestinationID  *uuid.UUID
	DestinationSet bool
	ProgramSlug    *string
	Slug           *string
	Name           *string
	Body           *string
	ComposeMode    *string
	AIPromptHint   *string
	Bindings       map[string]string
	BindingsSet    bool
	Active         *bool
}

func (s *Service) ListTemplates(ctx context.Context, f TemplateFilter) ([]models.MessageTemplate, error) {
	rows, err := s.repo.ListTemplates(ctx, repository.TemplateFilter(f))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load templates.", err)
	}
	return rows, nil
}

func (s *Service) GetTemplate(ctx context.Context, id uuid.UUID) (*models.MessageTemplate, error) {
	row, err := s.repo.FindTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Template not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load template.", err)
	}
	return row, nil
}

func (s *Service) CreateTemplate(ctx context.Context, in CreateTemplateInput) (*models.MessageTemplate, error) {
	slug := strings.TrimSpace(in.Slug)
	name := strings.TrimSpace(in.Name)
	body := strings.TrimSpace(in.Body)
	if slug == "" || name == "" || body == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Slug, name and body are required.")
	}
	row := &models.MessageTemplate{
		DestinationID: in.DestinationID,
		ProgramSlug:   strings.TrimSpace(in.ProgramSlug),
		Slug:          slug,
		Name:          name,
		Body:          body,
		ComposeMode:   normalizeComposeMode(in.ComposeMode),
		AIPromptHint:  strings.TrimSpace(in.AIPromptHint),
		Bindings:      jsonMapString(in.Bindings),
		Active:        in.Active,
	}
	if err := s.repo.CreateTemplate(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create template.", err)
	}
	return row, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, id uuid.UUID, in UpdateTemplateInput) (*models.MessageTemplate, error) {
	row, err := s.GetTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.DestinationSet {
		row.DestinationID = in.DestinationID
	}
	if in.ProgramSlug != nil {
		row.ProgramSlug = strings.TrimSpace(*in.ProgramSlug)
	}
	if in.Slug != nil {
		slug := strings.TrimSpace(*in.Slug)
		if slug == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Slug cannot be empty.")
		}
		row.Slug = slug
	}
	if in.Name != nil {
		row.Name = strings.TrimSpace(*in.Name)
	}
	if in.Body != nil {
		row.Body = strings.TrimSpace(*in.Body)
	}
	if in.ComposeMode != nil {
		row.ComposeMode = normalizeComposeMode(*in.ComposeMode)
	}
	if in.AIPromptHint != nil {
		row.AIPromptHint = strings.TrimSpace(*in.AIPromptHint)
	}
	if in.BindingsSet {
		row.Bindings = jsonMapString(in.Bindings)
	}
	if in.Active != nil {
		row.Active = *in.Active
	}
	if err := s.repo.SaveTemplate(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update template.", err)
	}
	return row, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteTemplate(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Template not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete template.", err)
	}
	return nil
}

type CreateJobInput struct {
	Name          string
	DestinationID uuid.UUID
	TemplateSlug  string
	ProgramAction string
	DataSource    map[string]any
	CronExpr      string
	Timezone      string
	Enabled       bool
}

type UpdateJobInput struct {
	Name          *string
	DestinationID *uuid.UUID
	TemplateSlug  *string
	ProgramAction *string
	DataSource    map[string]any
	DataSourceSet bool
	CronExpr      *string
	Timezone      *string
	Enabled       *bool
}

func (s *Service) ListJobs(ctx context.Context, enabledOnly bool) ([]models.ScheduledJob, error) {
	rows, err := s.repo.ListJobs(ctx, enabledOnly)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load jobs.", err)
	}
	return rows, nil
}

func (s *Service) ListDueJobs(ctx context.Context, before time.Time) ([]models.ScheduledJob, error) {
	if before.IsZero() {
		before = time.Now().UTC()
	}
	rows, err := s.repo.ListDueJobs(ctx, before)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load due jobs.", err)
	}
	return rows, nil
}

func (s *Service) GetJob(ctx context.Context, id uuid.UUID) (*models.ScheduledJob, error) {
	row, err := s.repo.FindJob(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Job not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load job.", err)
	}
	return row, nil
}

func (s *Service) CreateJob(ctx context.Context, in CreateJobInput) (*models.ScheduledJob, error) {
	name := strings.TrimSpace(in.Name)
	cronExpr := strings.TrimSpace(in.CronExpr)
	if name == "" || cronExpr == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Name and cron expression are required.")
	}
	if _, err := s.GetDestination(ctx, in.DestinationID); err != nil {
		return nil, err
	}
	tz := strings.TrimSpace(in.Timezone)
	if tz == "" {
		tz = "America/Sao_Paulo"
	}
	next, err := computeNextRun(cronExpr, tz, time.Now().UTC())
	if err != nil {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Invalid cron expression.")
	}
	row := &models.ScheduledJob{
		Name:          name,
		DestinationID: in.DestinationID,
		TemplateSlug:  strings.TrimSpace(in.TemplateSlug),
		ProgramAction: strings.TrimSpace(in.ProgramAction),
		DataSource:    jsonMap(in.DataSource),
		CronExpr:      cronExpr,
		Timezone:      tz,
		Enabled:       in.Enabled,
		NextRunAt:     &next,
	}
	if err := s.repo.CreateJob(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create job.", err)
	}
	return row, nil
}

func (s *Service) UpdateJob(ctx context.Context, id uuid.UUID, in UpdateJobInput) (*models.ScheduledJob, error) {
	row, err := s.GetJob(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		row.Name = strings.TrimSpace(*in.Name)
	}
	if in.DestinationID != nil {
		if _, err := s.GetDestination(ctx, *in.DestinationID); err != nil {
			return nil, err
		}
		row.DestinationID = *in.DestinationID
	}
	if in.TemplateSlug != nil {
		row.TemplateSlug = strings.TrimSpace(*in.TemplateSlug)
	}
	if in.ProgramAction != nil {
		row.ProgramAction = strings.TrimSpace(*in.ProgramAction)
	}
	if in.DataSourceSet {
		row.DataSource = jsonMap(in.DataSource)
	}
	recalc := false
	if in.CronExpr != nil {
		row.CronExpr = strings.TrimSpace(*in.CronExpr)
		recalc = true
	}
	if in.Timezone != nil {
		row.Timezone = strings.TrimSpace(*in.Timezone)
		recalc = true
	}
	if in.Enabled != nil {
		row.Enabled = *in.Enabled
	}
	if recalc {
		next, err := computeNextRun(row.CronExpr, row.Timezone, time.Now().UTC())
		if err != nil {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Invalid cron expression.")
		}
		row.NextRunAt = &next
	}
	if err := s.repo.SaveJob(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update job.", err)
	}
	return row, nil
}

func (s *Service) DeleteJob(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteJob(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Job not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete job.", err)
	}
	return nil
}

func (s *Service) RecordDelivery(ctx context.Context, d *models.MessageDelivery) error {
	if err := s.repo.CreateDelivery(ctx, d); err != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to record delivery.", err)
	}
	return nil
}

func (s *Service) ListDeliveries(ctx context.Context, destinationID *uuid.UUID, limit int) ([]models.MessageDelivery, error) {
	rows, err := s.repo.ListDeliveries(ctx, destinationID, limit)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load deliveries.", err)
	}
	return rows, nil
}

func (s *Service) FindTemplateBySlug(ctx context.Context, slug string, destinationID uuid.UUID) (*models.MessageTemplate, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, apperrors.NotFound(apperrors.CodeInternal, "Template not found.")
	}
	destID := &destinationID
	row, err := s.repo.FindTemplateBySlug(ctx, "", slug, destID)
	if err == nil {
		return row, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load template.", err)
	}
	row, err = s.repo.FindTemplateBySlug(ctx, "", slug, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Template not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load template.", err)
	}
	return row, nil
}

func (s *Service) TemplateBodyForProgram(ctx context.Context, programSlug, slug string) (string, bool) {
	row, err := s.repo.FindTemplateBySlug(ctx, strings.TrimSpace(programSlug), strings.TrimSpace(slug), nil)
	if err != nil {
		return "", false
	}
	return row.Body, true
}

func normalizeChannel(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case models.ChannelWhatsApp, models.ChannelTelegram:
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return ""
	}
}

func normalizeComposeMode(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case models.ComposeModeAIAssisted:
		return models.ComposeModeAIAssisted
	default:
		return models.ComposeModeStatic
	}
}

func jsonSlice(items []string) datatypes.JSON {
	if len(items) == 0 {
		return datatypes.JSON([]byte("[]"))
	}
	b, _ := json.Marshal(items)
	return datatypes.JSON(b)
}

func jsonMap(m map[string]any) datatypes.JSON {
	if len(m) == 0 {
		return datatypes.JSON([]byte("{}"))
	}
	b, _ := json.Marshal(m)
	return datatypes.JSON(b)
}

func jsonMapString(m map[string]string) datatypes.JSON {
	if len(m) == 0 {
		return datatypes.JSON([]byte("{}"))
	}
	b, _ := json.Marshal(m)
	return datatypes.JSON(b)
}

func computeNextRun(cronExpr, timezone string, from time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := parser.Parse(cronExpr)
	if err != nil {
		return time.Time{}, err
	}
	return sched.Next(from.In(loc)).UTC(), nil
}
