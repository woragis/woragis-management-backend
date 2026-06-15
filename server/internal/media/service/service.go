package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/media/repository"
	"github.com/woragis/management/backend/server/internal/media/storage"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

const maxUploadBytes = 10 << 20 // 10 MiB

type Service struct {
	repo    *repository.Repository
	store   *storage.Local
	baseURL string
}

func New(repo *repository.Repository, store *storage.Local, baseURL string) *Service {
	return &Service{repo: repo, store: store, baseURL: strings.TrimRight(baseURL, "/")}
}

type UploadInput struct {
	Filename string
	MimeType string
	AltText  string
	Reader   io.Reader
}

func (s *Service) List(ctx context.Context) ([]models.MediaAsset, error) {
	out, err := s.repo.List(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMediaListV1ServiceLoadFailed, apperrors.MsgMediaListV1ServiceLoadFailed, err)
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.MediaAsset, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeMediaGetV1ServiceNotFound, apperrors.MsgMediaGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeMediaListV1ServiceLoadFailed, apperrors.MsgMediaListV1ServiceLoadFailed, err)
	}
	return m, nil
}

func (s *Service) Upload(ctx context.Context, in UploadInput) (*models.MediaAsset, error) {
	if in.Reader == nil {
		return nil, apperrors.Invalid(apperrors.CodeMediaPostV1ServiceFileMissing, apperrors.MsgMediaPostV1ServiceFileMissing)
	}
	limited := io.LimitReader(in.Reader, maxUploadBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMediaPostV1ServiceCreateFailed, apperrors.MsgMediaPostV1ServiceCreateFailed, err)
	}
	if len(data) == 0 {
		return nil, apperrors.Invalid(apperrors.CodeMediaPostV1ServiceFileMissing, apperrors.MsgMediaPostV1ServiceFileMissing)
	}
	if len(data) > maxUploadBytes {
		return nil, apperrors.Invalid(apperrors.CodeMediaPostV1ServiceFileTooLarge, apperrors.MsgMediaPostV1ServiceFileTooLarge)
	}

	id := uuid.New()
	ext := filepath.Ext(in.Filename)
	if ext == "" {
		ext = extFromMime(in.MimeType)
	}
	key := id.String() + ext
	if _, err := s.store.Save(key, bytes.NewReader(data)); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMediaPostV1ServiceCreateFailed, apperrors.MsgMediaPostV1ServiceCreateFailed, err)
	}

	mimeType := strings.TrimSpace(in.MimeType)
	if mimeType == "" {
		mimeType = mime.TypeByExtension(ext)
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	asset := &models.MediaAsset{
		ID:         id,
		Filename:   strings.TrimSpace(in.Filename),
		MimeType:   mimeType,
		SizeBytes:  int64(len(data)),
		StorageKey: key,
		PublicURL:  fmt.Sprintf("%s/%s/file", s.baseURL, id.String()),
		AltText:    strings.TrimSpace(in.AltText),
	}
	if err := s.repo.Create(ctx, asset); err != nil {
		_ = s.store.Delete(key)
		return nil, apperrors.InternalCause(apperrors.CodeMediaPostV1ServiceCreateFailed, apperrors.MsgMediaPostV1ServiceCreateFailed, err)
	}
	return asset, nil
}

func (s *Service) Open(ctx context.Context, id uuid.UUID) (*models.MediaAsset, io.ReadCloser, error) {
	m, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	f, err := s.store.Open(m.StorageKey)
	if err != nil {
		return nil, nil, apperrors.NotFound(apperrors.CodeMediaGetV1ServiceNotFound, apperrors.MsgMediaGetV1ServiceNotFound)
	}
	return m, f, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	m, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeMediaDeleteV1ServiceNotFound, apperrors.MsgMediaDeleteV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeMediaPostV1ServiceCreateFailed, apperrors.MsgMediaPostV1ServiceCreateFailed, err)
	}
	_ = s.store.Delete(m.StorageKey)
	return nil
}

func extFromMime(mimeType string) string {
	switch strings.ToLower(mimeType) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "application/pdf":
		return ".pdf"
	default:
		return ".bin"
	}
}
