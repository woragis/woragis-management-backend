package service

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var templateVarRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

type TemplateVars map[string]string

type RenderedTemplate struct {
	TemplateSlug string `json:"templateSlug"`
	TemplateName string `json:"templateName"`
	Platform     string `json:"platform"`
	Goal         string `json:"goal"`
	Body         string `json:"body"`
	CharLimit    int    `json:"charLimit"`
	CharCount    int    `json:"charCount"`
	OverLimit    bool   `json:"overLimit"`
}

func TemplateVarsFromProject(p *models.Project) TemplateVars {
	stack := parseStackSlice(p.Stack)
	return TemplateVars{
		"projectName":      p.Name,
		"projectSlug":      p.Slug,
		"shortDescription": p.ShortDescription,
		"demoUrl":          p.DemoURL,
		"githubUrl":        p.GithubURL,
		"repoUrl":          p.RepoURL,
		"stack":            strings.Join(stack, ", "),
	}
}

func ApplyTemplateBody(body string, vars TemplateVars) string {
	return templateVarRe.ReplaceAllStringFunc(body, func(match string) string {
		key := strings.Trim(match, "{}")
		if v, ok := vars[key]; ok {
			return v
		}
		return ""
	})
}

func platformCharLimit(platform string) int {
	switch platform {
	case models.SocialPlatformTwitter:
		return 280
	case models.SocialPlatformLinkedIn:
		return 3000
	case models.SocialPlatformReddit:
		return 40000
	default:
		return 3000
	}
}

func (s *Service) RenderTemplate(ctx context.Context, slug string, vars TemplateVars) (*RenderedTemplate, error) {
	slug = strings.TrimSpace(strings.ToLower(slug))
	if slug == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Template slug is required.")
	}
	tpl, err := s.repo.FindTemplateBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Template not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load template.", err)
	}
	if !tpl.Active {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Template is inactive.")
	}
	body := ApplyTemplateBody(tpl.Body, vars)
	platform := tpl.Platform
	if platform == "any" {
		platform = models.SocialPlatformLinkedIn
	}
	limit := platformCharLimit(platform)
	count := len([]rune(body))
	return &RenderedTemplate{
		TemplateSlug: tpl.Slug,
		TemplateName: tpl.Name,
		Platform:     platform,
		Goal:         tpl.Goal,
		Body:         body,
		CharLimit:    limit,
		CharCount:    count,
		OverLimit:    count > limit,
	}, nil
}

func parseStackSlice(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return nil
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}
