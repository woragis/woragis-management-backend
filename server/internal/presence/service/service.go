package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	"github.com/woragis/management/backend/server/internal/presence/repository"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

type CampaignFilter struct {
	Goal       string
	ProjectID  *uuid.UUID
	ActiveOnly bool
}

type CreateCampaignInput struct {
	Name        string
	Goal        string
	Description string
	ProjectID   *uuid.UUID
	StartDate   *time.Time
	EndDate     *time.Time
	Active      bool
}

type UpdateCampaignInput struct {
	Name        *string
	Goal        *string
	Description *string
	ProjectID   *uuid.UUID
	ProjectSet  bool
	StartDate   *time.Time
	StartSet    bool
	EndDate     *time.Time
	EndSet      bool
	Active      *bool
}

type TemplateFilter struct {
	Platform   string
	Goal       string
	ActiveOnly bool
}

type CreateTemplateInput struct {
	Slug     string
	Name     string
	Platform string
	Goal     string
	Body     string
	Active   bool
}

type UpdateTemplateInput struct {
	Slug     *string
	Name     *string
	Platform *string
	Goal     *string
	Body     *string
	Active   *bool
}

type PostFilter struct {
	Platform   string
	Goal       string
	Status     string
	ProjectID  *uuid.UUID
	CampaignID *uuid.UUID
	Limit      int
}

type CreatePostInput struct {
	ProjectID    *uuid.UUID
	CampaignID   *uuid.UUID
	Platform     string
	Goal         string
	Status       string
	Title        string
	Body         string
	Hook         string
	CTA          string
	TemplateSlug string
	ScheduledAt  *time.Time
	PublishedAt  *time.Time
	PublishedURL string
	Notes        string
}

type UpdatePostInput struct {
	ProjectID    *uuid.UUID
	ProjectSet   bool
	CampaignID   *uuid.UUID
	CampaignSet  bool
	Platform     *string
	Goal         *string
	Status       *string
	Title        *string
	Body         *string
	Hook         *string
	CTA          *string
	TemplateSlug *string
	ScheduledAt  *time.Time
	ScheduledSet bool
	PublishedAt  *time.Time
	PublishedSet bool
	PublishedURL *string
	Notes        *string
}

func (s *Service) ListCampaigns(ctx context.Context, f CampaignFilter) ([]models.SocialCampaign, error) {
	rows, err := s.repo.ListCampaigns(ctx, repository.CampaignFilter{
		Goal:       normalizeGoal(f.Goal),
		ProjectID:  f.ProjectID,
		ActiveOnly: f.ActiveOnly,
	})
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load campaigns.", err)
	}
	return rows, nil
}

func (s *Service) GetCampaign(ctx context.Context, id uuid.UUID) (*models.SocialCampaign, error) {
	row, err := s.repo.FindCampaign(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Campaign not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load campaign.", err)
	}
	return row, nil
}

func (s *Service) CreateCampaign(ctx context.Context, in CreateCampaignInput) (*models.SocialCampaign, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Name is required.")
	}
	row := &models.SocialCampaign{
		Name:        name,
		Goal:        normalizeGoal(in.Goal),
		Description: strings.TrimSpace(in.Description),
		ProjectID:   in.ProjectID,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		Active:      in.Active,
	}
	if err := s.repo.CreateCampaign(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create campaign.", err)
	}
	return row, nil
}

func (s *Service) UpdateCampaign(ctx context.Context, id uuid.UUID, in UpdateCampaignInput) (*models.SocialCampaign, error) {
	row, err := s.GetCampaign(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Name cannot be empty.")
		}
		row.Name = name
	}
	if in.Goal != nil {
		row.Goal = normalizeGoal(*in.Goal)
	}
	if in.Description != nil {
		row.Description = strings.TrimSpace(*in.Description)
	}
	if in.ProjectSet {
		row.ProjectID = in.ProjectID
	}
	if in.StartSet {
		row.StartDate = in.StartDate
	}
	if in.EndSet {
		row.EndDate = in.EndDate
	}
	if in.Active != nil {
		row.Active = *in.Active
	}
	if err := s.repo.SaveCampaign(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update campaign.", err)
	}
	return row, nil
}

func (s *Service) DeleteCampaign(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteCampaign(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Campaign not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete campaign.", err)
	}
	return nil
}

func (s *Service) ListTemplates(ctx context.Context, f TemplateFilter) ([]models.PostTemplate, error) {
	platform := strings.TrimSpace(strings.ToLower(f.Platform))
	goal := strings.TrimSpace(strings.ToLower(f.Goal))
	rows, err := s.repo.ListTemplates(ctx, repository.TemplateFilter{
		Platform:   platform,
		Goal:       goal,
		ActiveOnly: f.ActiveOnly,
	})
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load templates.", err)
	}
	return rows, nil
}

func (s *Service) GetTemplate(ctx context.Context, id uuid.UUID) (*models.PostTemplate, error) {
	row, err := s.repo.FindTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Template not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load template.", err)
	}
	return row, nil
}

func (s *Service) CreateTemplate(ctx context.Context, in CreateTemplateInput) (*models.PostTemplate, error) {
	slug := strings.TrimSpace(strings.ToLower(in.Slug))
	name := strings.TrimSpace(in.Name)
	body := strings.TrimSpace(in.Body)
	if slug == "" || name == "" || body == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Slug, name, and body are required.")
	}
	row := &models.PostTemplate{
		Slug:     slug,
		Name:     name,
		Platform: normalizePlatform(in.Platform),
		Goal:     normalizeGoal(in.Goal),
		Body:     body,
		Active:   in.Active,
	}
	if err := s.repo.CreateTemplate(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create template.", err)
	}
	return row, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, id uuid.UUID, in UpdateTemplateInput) (*models.PostTemplate, error) {
	row, err := s.GetTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Slug != nil {
		slug := strings.TrimSpace(strings.ToLower(*in.Slug))
		if slug == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Slug cannot be empty.")
		}
		row.Slug = slug
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Name cannot be empty.")
		}
		row.Name = name
	}
	if in.Platform != nil {
		row.Platform = normalizePlatform(*in.Platform)
	}
	if in.Goal != nil {
		row.Goal = normalizeGoal(*in.Goal)
	}
	if in.Body != nil {
		body := strings.TrimSpace(*in.Body)
		if body == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Body cannot be empty.")
		}
		row.Body = body
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

func (s *Service) ListPosts(ctx context.Context, f PostFilter) ([]models.SocialPost, error) {
	platform := strings.TrimSpace(strings.ToLower(f.Platform))
	goal := strings.TrimSpace(strings.ToLower(f.Goal))
	status := strings.TrimSpace(strings.ToLower(f.Status))
	rows, err := s.repo.ListPosts(ctx, repository.PostFilter{
		Platform:   platform,
		Goal:       goal,
		Status:     status,
		ProjectID:  f.ProjectID,
		CampaignID: f.CampaignID,
		Limit:      f.Limit,
	})
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load posts.", err)
	}
	return rows, nil
}

func (s *Service) GetPost(ctx context.Context, id uuid.UUID) (*models.SocialPost, error) {
	row, err := s.repo.FindPost(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Post not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load post.", err)
	}
	return row, nil
}

func (s *Service) CreatePost(ctx context.Context, in CreatePostInput) (*models.SocialPost, error) {
	body := strings.TrimSpace(in.Body)
	if body == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Body is required.")
	}
	platform := normalizePlatform(in.Platform)
	if platform == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Platform must be linkedin, reddit, or twitter.")
	}
	row := &models.SocialPost{
		ProjectID:    in.ProjectID,
		CampaignID:   in.CampaignID,
		Platform:     platform,
		Goal:         normalizeGoal(in.Goal),
		Status:       normalizePostStatus(in.Status),
		Title:        strings.TrimSpace(in.Title),
		Body:         body,
		Hook:         strings.TrimSpace(in.Hook),
		CTA:          strings.TrimSpace(in.CTA),
		TemplateSlug: strings.TrimSpace(in.TemplateSlug),
		ScheduledAt:  in.ScheduledAt,
		PublishedAt:  in.PublishedAt,
		PublishedURL: strings.TrimSpace(in.PublishedURL),
		Notes:        strings.TrimSpace(in.Notes),
	}
	if row.Status == models.SocialPostStatusPublished && row.PublishedAt == nil {
		now := time.Now().UTC()
		row.PublishedAt = &now
	}
	if err := s.repo.CreatePost(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create post.", err)
	}
	return row, nil
}

func (s *Service) UpdatePost(ctx context.Context, id uuid.UUID, in UpdatePostInput) (*models.SocialPost, error) {
	row, err := s.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.ProjectSet {
		row.ProjectID = in.ProjectID
	}
	if in.CampaignSet {
		row.CampaignID = in.CampaignID
	}
	if in.Platform != nil {
		platform := normalizePlatform(*in.Platform)
		if platform == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Invalid platform.")
		}
		row.Platform = platform
	}
	if in.Goal != nil {
		row.Goal = normalizeGoal(*in.Goal)
	}
	if in.Status != nil {
		row.Status = normalizePostStatus(*in.Status)
		if row.Status == models.SocialPostStatusPublished && row.PublishedAt == nil {
			now := time.Now().UTC()
			row.PublishedAt = &now
		}
	}
	if in.Title != nil {
		row.Title = strings.TrimSpace(*in.Title)
	}
	if in.Body != nil {
		body := strings.TrimSpace(*in.Body)
		if body == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Body cannot be empty.")
		}
		row.Body = body
	}
	if in.Hook != nil {
		row.Hook = strings.TrimSpace(*in.Hook)
	}
	if in.CTA != nil {
		row.CTA = strings.TrimSpace(*in.CTA)
	}
	if in.TemplateSlug != nil {
		row.TemplateSlug = strings.TrimSpace(*in.TemplateSlug)
	}
	if in.ScheduledSet {
		row.ScheduledAt = in.ScheduledAt
	}
	if in.PublishedSet {
		row.PublishedAt = in.PublishedAt
	}
	if in.PublishedURL != nil {
		row.PublishedURL = strings.TrimSpace(*in.PublishedURL)
	}
	if in.Notes != nil {
		row.Notes = strings.TrimSpace(*in.Notes)
	}
	if err := s.repo.SavePost(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update post.", err)
	}
	return row, nil
}

func (s *Service) DeletePost(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeletePost(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Post not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete post.", err)
	}
	return nil
}

func normalizePlatform(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case models.SocialPlatformLinkedIn, models.SocialPlatformReddit, models.SocialPlatformTwitter:
		return strings.TrimSpace(strings.ToLower(v))
	case "any", "":
		return "any"
	default:
		return ""
	}
}

func normalizeGoal(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case models.SocialGoalJobHunting, models.SocialGoalRevenue, models.SocialGoalLaunch,
		models.SocialGoalVisibility, models.SocialGoalAcademic, models.SocialGoalCommunity:
		return strings.TrimSpace(strings.ToLower(v))
	default:
		return models.SocialGoalVisibility
	}
}

func normalizePostStatus(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case models.SocialPostStatusDraft, models.SocialPostStatusScheduled,
		models.SocialPostStatusPublished, models.SocialPostStatusCancelled:
		return strings.TrimSpace(strings.ToLower(v))
	default:
		return models.SocialPostStatusDraft
	}
}
