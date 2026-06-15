package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	"github.com/woragis/management/backend/server/internal/profile/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

type UpdateInput struct {
	DisplayName   *string
	Headline      *string
	Bio           *string
	AvatarID      *uuid.UUID
	AvatarSet     bool
	Location      *string
	Availability  *string
	ResumeAssetID *uuid.UUID
	ResumeSet     bool
	SocialLinks   []models.SocialLink
	SocialSet     bool
}

func (s *Service) Get(ctx context.Context) (*models.Profile, error) {
	p, err := s.repo.EnsureDefault(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProfileGetV1ServiceLoadFailed, apperrors.MsgProfileGetV1ServiceLoadFailed, err)
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, in UpdateInput) (*models.Profile, error) {
	p, err := s.Get(ctx)
	if err != nil {
		return nil, err
	}
	if in.DisplayName != nil {
		name := strings.TrimSpace(*in.DisplayName)
		if name == "" {
			return nil, apperrors.Invalid(apperrors.CodeProfilePatchV1ServiceDisplayNameEmpty, apperrors.MsgProfilePatchV1ServiceDisplayNameEmpty)
		}
		p.DisplayName = name
	}
	if in.Headline != nil {
		p.Headline = strings.TrimSpace(*in.Headline)
	}
	if in.Bio != nil {
		p.Bio = strings.TrimSpace(*in.Bio)
	}
	if in.AvatarSet {
		p.AvatarID = in.AvatarID
	}
	if in.Location != nil {
		p.Location = strings.TrimSpace(*in.Location)
	}
	if in.Availability != nil {
		p.Availability = normalizeAvailability(*in.Availability)
	}
	if in.ResumeSet {
		p.ResumeAssetID = in.ResumeAssetID
	}
	if in.SocialSet {
		p.SocialLinks = socialJSON(in.SocialLinks)
	}
	if err := s.repo.Save(ctx, p); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProfilePatchV1ServiceUpdateFailed, apperrors.MsgProfilePatchV1ServiceUpdateFailed, err)
	}
	return s.Get(ctx)
}

func (s *Service) GetPublic(ctx context.Context) (*models.Profile, error) {
	p, err := s.repo.GetDefault(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeProfileGetV1ServiceNotFound, apperrors.MsgProfileGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeProfileGetV1ServiceLoadFailed, apperrors.MsgProfileGetV1ServiceLoadFailed, err)
	}
	return p, nil
}

func normalizeAvailability(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "open_to_work", "freelancing":
		return strings.TrimSpace(strings.ToLower(v))
	default:
		return "not_available"
	}
}

func socialJSON(links []models.SocialLink) datatypes.JSON {
	if links == nil {
		return datatypes.JSON([]byte("[]"))
	}
	b, _ := json.Marshal(links)
	return datatypes.JSON(b)
}
