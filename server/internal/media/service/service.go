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

type Service struct {
	repo    *repository.Repository
	store   storage.BlobStore
	baseURL string
}

func New(repo *repository.Repository, store storage.BlobStore, baseURL string) *Service {
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
	mimeType := strings.TrimSpace(in.MimeType)
	limit := int64(maxBytesForMime(mimeType))
	limited := io.LimitReader(in.Reader, limit+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMediaPostV1ServiceCreateFailed, apperrors.MsgMediaPostV1ServiceCreateFailed, err)
	}
	if len(data) == 0 {
		return nil, apperrors.Invalid(apperrors.CodeMediaPostV1ServiceFileMissing, apperrors.MsgMediaPostV1ServiceFileMissing)
	}
	if int64(len(data)) > limit {
		return nil, apperrors.Invalid(apperrors.CodeMediaPostV1ServiceFileTooLarge, apperrors.MsgMediaPostV1ServiceFileTooLarge)
	}

	id := uuid.New()
	ext := filepath.Ext(in.Filename)
	if ext == "" {
		ext = extFromMime(mimeType)
	}
	key := id.String() + ext

	if mimeType == "" {
		mimeType = mime.TypeByExtension(ext)
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if !IsGalleryMime(mimeType) && mimeType != "application/pdf" {
		return nil, apperrors.Invalid(apperrors.CodeMediaPostV1ServiceFileMissing, "Unsupported file type.")
	}

	if _, err := s.store.Save(ctx, key, bytes.NewReader(data), mimeType); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMediaPostV1ServiceCreateFailed, apperrors.MsgMediaPostV1ServiceCreateFailed, err)
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
		_ = s.store.Delete(ctx, key)
		return nil, apperrors.InternalCause(apperrors.CodeMediaPostV1ServiceCreateFailed, apperrors.MsgMediaPostV1ServiceCreateFailed, err)
	}
	return asset, nil
}

func (s *Service) Open(ctx context.Context, id uuid.UUID) (*models.MediaAsset, io.ReadCloser, error) {
	m, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	f, err := s.store.Open(ctx, m.StorageKey)
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
	_ = s.store.Delete(ctx, m.StorageKey)
	return nil
}
