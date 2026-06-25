package service

import (
	"strings"

	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func normalizeAccessLevel(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "public", "secret":
		return strings.TrimSpace(strings.ToLower(v))
	default:
		return "private"
	}
}

func applyAccessLevel(p *models.Project, level string) {
	level = normalizeAccessLevel(level)
	p.AccessLevel = level
	switch level {
	case "public":
		p.IsPublic = true
	case "secret", "private":
		p.IsPublic = false
	}
	if level == "secret" {
		p.Featured = false
	}
}

func (s *Service) verifySecretUnlock(password string) error {
	if len(s.secretUnlockHash) == 0 {
		return apperrors.Invalid(apperrors.CodeProjectSecretUnlockUnavailable, apperrors.MsgProjectSecretUnlockUnavailable)
	}
	if strings.TrimSpace(password) == "" {
		return apperrors.Invalid(apperrors.CodeProjectSecretUnlockRequired, apperrors.MsgProjectSecretUnlockRequired)
	}
	if err := bcrypt.CompareHashAndPassword(s.secretUnlockHash, []byte(password)); err != nil {
		return apperrors.Invalid(apperrors.CodeProjectSecretUnlockInvalid, apperrors.MsgProjectSecretUnlockInvalid)
	}
	return nil
}

func (s *Service) guardSecretProjectChange(p *models.Project, in UpdateProjectInput) error {
	if p.AccessLevel != "secret" {
		return nil
	}
	needsUnlock := false
	if in.AccessLevel != nil && normalizeAccessLevel(*in.AccessLevel) != "secret" {
		needsUnlock = true
	}
	if in.IsPublic != nil && *in.IsPublic {
		needsUnlock = true
	}
	if in.Featured != nil && *in.Featured {
		needsUnlock = true
	}
	if !needsUnlock {
		return nil
	}
	return s.verifySecretUnlock(in.SecretUnlockPassword)
}
