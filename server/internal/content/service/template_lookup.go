package service

import (
	"context"
)

func (s *Service) templateBody(ctx context.Context, programSlug, slug string) (string, error) {
	if s.msgTemplates != nil {
		if body, ok := s.msgTemplates.TemplateBodyForProgram(ctx, programSlug, slug); ok && body != "" {
			return body, nil
		}
	}
	tpl, err := s.repo.GetWhatsappTemplateBySlug(ctx, slug)
	if err != nil {
		return "", err
	}
	return tpl.Body, nil
}
