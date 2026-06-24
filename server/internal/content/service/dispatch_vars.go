package service

import (
	"context"
	"strings"
	"time"

	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/content/templaterender"
	"github.com/woragis/management/backend/server/internal/models"
)

type DispatchVarsResult struct {
	Vars         templaterender.Vars
	VideoID      string
	TemplateSlug string
	Skip         bool
	SkipReason   string
}

func (s *Service) ResolveDispatchVars(ctx context.Context, dispatchType, dateStr string) (*DispatchVarsResult, error) {
	settings, err := s.repo.EnsureSettings(ctx)
	if err != nil {
		return nil, err
	}
	day := templaterender.TodayInTZ(settings.Timezone)
	if dateStr != "" {
		day, err = templaterender.ParseDateInTZ(dateStr, settings.Timezone)
		if err != nil {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Invalid date.")
		}
	}

	switch dispatchType {
	case "problem":
		return s.varsProblem(ctx, settings, day)
	case "discussion":
		return s.varsDiscussion(ctx, settings, day)
	case "solution":
		return s.varsSolution(ctx, settings, day)
	case "weekly":
		return s.varsWeekly(ctx, settings, day)
	default:
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Invalid dispatch type.")
	}
}

func (s *Service) varsProblem(ctx context.Context, settings *models.LeetcodeChannelSettings, day time.Time) (*DispatchVarsResult, error) {
	video, err := s.repo.GetVideoByProblemDate(ctx, day)
	if err != nil {
		return &DispatchVarsResult{Skip: true, SkipReason: "no video for date", TemplateSlug: models.WhatsappTplProblemDaily}, nil
	}
	if video.WhatsappProblemSentAt != nil {
		return &DispatchVarsResult{Skip: true, SkipReason: "already sent", TemplateSlug: models.WhatsappTplProblemDaily, VideoID: video.ID.String()}, nil
	}
	return &DispatchVarsResult{
		Vars:         templaterender.FromVideo(video, settings),
		VideoID:      video.ID.String(),
		TemplateSlug: models.WhatsappTplProblemDaily,
	}, nil
}

func (s *Service) varsDiscussion(ctx context.Context, settings *models.LeetcodeChannelSettings, day time.Time) (*DispatchVarsResult, error) {
	if !settings.DiscussionEnabled {
		return &DispatchVarsResult{Skip: true, SkipReason: "discussion disabled", TemplateSlug: models.WhatsappTplDiscussionNudge}, nil
	}
	video, err := s.repo.GetVideoByProblemDate(ctx, day)
	if err != nil {
		return &DispatchVarsResult{Skip: true, SkipReason: "no video for date", TemplateSlug: models.WhatsappTplDiscussionNudge}, nil
	}
	if video.WhatsappDiscussionSentAt != nil {
		return &DispatchVarsResult{Skip: true, SkipReason: "already sent", TemplateSlug: models.WhatsappTplDiscussionNudge, VideoID: video.ID.String()}, nil
	}
	return &DispatchVarsResult{
		Vars:         templaterender.FromVideo(video, settings),
		VideoID:      video.ID.String(),
		TemplateSlug: models.WhatsappTplDiscussionNudge,
	}, nil
}

func (s *Service) varsSolution(ctx context.Context, settings *models.LeetcodeChannelSettings, day time.Time) (*DispatchVarsResult, error) {
	video, err := s.repo.GetVideoByProblemDate(ctx, day)
	if err != nil {
		return &DispatchVarsResult{Skip: true, SkipReason: "no video for date", TemplateSlug: models.WhatsappTplSolutionVideo}, nil
	}
	if video.WhatsappSolutionSentAt != nil {
		return &DispatchVarsResult{Skip: true, SkipReason: "already sent", TemplateSlug: models.WhatsappTplSolutionVideo, VideoID: video.ID.String()}, nil
	}
	if video.YoutubeURL == nil || strings.TrimSpace(*video.YoutubeURL) == "" {
		return &DispatchVarsResult{Skip: true, SkipReason: "youtube url missing", TemplateSlug: models.WhatsappTplSolutionVideo, VideoID: video.ID.String()}, nil
	}
	return &DispatchVarsResult{
		Vars:         templaterender.FromVideo(video, settings),
		VideoID:      video.ID.String(),
		TemplateSlug: models.WhatsappTplSolutionVideo,
	}, nil
}

func (s *Service) varsWeekly(ctx context.Context, settings *models.LeetcodeChannelSettings, day time.Time) (*DispatchVarsResult, error) {
	loc, _ := time.LoadLocation(settings.Timezone)
	if loc == nil {
		loc = time.UTC
	}
	day = day.In(loc)
	weekday := int(day.Weekday())
	start := day.AddDate(0, 0, -weekday)
	end := start.AddDate(0, 0, 7)
	videos, err := s.repo.ListVideosWithProblemSentBetween(ctx, start, end)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return &DispatchVarsResult{Skip: true, SkipReason: "no problems this week", TemplateSlug: models.WhatsappTplWeeklySummary}, nil
	}
	vars := templaterender.Vars{ProblemList: templaterender.FormatProblemList(videos)}
	if settings.NextTheme != nil {
		vars.NextTheme = strings.TrimSpace(*settings.NextTheme)
	}
	return &DispatchVarsResult{Vars: vars, TemplateSlug: models.WhatsappTplWeeklySummary}, nil
}
