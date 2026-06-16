package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/devproject/repository"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

type ListFilter struct {
	Status   string
	IsPublic *bool
	Featured *bool
	Query    string
}

type Dashboard struct {
	ProjectCount        int64                  `json:"projectCount"`
	PublicProjectCount  int64                  `json:"publicProjectCount"`
	ActiveProjectCount  int64                  `json:"activeProjectCount"`
	MediaCount          int64                  `json:"mediaCount"`
	SecretsExpiringSoon []models.ProjectSecret `json:"secretsExpiringSoon"`
	DomainsExpiringSoon []models.ProjectDomain `json:"domainsExpiringSoon"`
}

type MediaCounter interface {
	Count(ctx context.Context) (int64, error)
}

func (s *Service) ListFiltered(ctx context.Context, f ListFilter) ([]models.Project, error) {
	out, err := s.repo.ListProjectsFiltered(ctx, repository.ListFilter{
		Status:   f.Status,
		IsPublic: f.IsPublic,
		Featured: f.Featured,
		Query:    f.Query,
	})
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	return out, nil
}

func (s *Service) Dashboard(ctx context.Context, media MediaCounter) (*Dashboard, error) {
	total, err := s.repo.CountProjects(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	public, err := s.repo.CountProjectsPublic(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	active, err := s.repo.CountProjectsActive(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	var mediaCount int64
	if media != nil {
		mediaCount, err = media.Count(ctx)
		if err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeMediaListV1ServiceLoadFailed, apperrors.MsgMediaListV1ServiceLoadFailed, err)
		}
	}
	threshold := time.Now().Add(30 * 24 * time.Hour)
	secrets, err := s.repo.ListSecretsExpiringBefore(ctx, threshold)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	domains, err := s.repo.ListDomainsExpiringBefore(ctx, threshold)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	return &Dashboard{
		ProjectCount:        total,
		PublicProjectCount:  public,
		ActiveProjectCount:  active,
		MediaCount:          mediaCount,
		SecretsExpiringSoon: secrets,
		DomainsExpiringSoon: domains,
	}, nil
}

type CreateEnvInput struct {
	Key         string
	Value       string
	Environment string
	Notes       string
}

func (s *Service) ListEnvs(ctx context.Context, projectID uuid.UUID) ([]models.ProjectEnv, error) {
	if _, err := s.GetByID(ctx, projectID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListEnvs(ctx, projectID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	return rows, nil
}

func (s *Service) CreateEnv(ctx context.Context, projectID uuid.UUID, in CreateEnvInput) (*models.ProjectEnv, error) {
	if _, err := s.GetByID(ctx, projectID); err != nil {
		return nil, err
	}
	key := strings.TrimSpace(in.Key)
	if key == "" {
		return nil, apperrors.Invalid(apperrors.CodeProjectEnvPostV1ServiceKeyEmpty, apperrors.MsgProjectEnvPostV1ServiceKeyEmpty)
	}
	e := &models.ProjectEnv{
		ProjectID:   projectID,
		Key:         key,
		Value:       in.Value,
		Environment: normalizeEnvironment(in.Environment),
		Notes:       strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateEnv(ctx, e); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectEnvPostV1ServiceCreateFailed, apperrors.MsgProjectEnvPostV1ServiceCreateFailed, err)
	}
	return e, nil
}

func (s *Service) DeleteEnv(ctx context.Context, projectID, envID uuid.UUID) error {
	if err := s.repo.DeleteEnv(ctx, projectID, envID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeProjectEnvDeleteV1ServiceNotFound, apperrors.MsgProjectEnvDeleteV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeProjectEnvPostV1ServiceCreateFailed, apperrors.MsgProjectEnvPostV1ServiceCreateFailed, err)
	}
	return nil
}
