package storage

import (
	"context"
	"io"
)

type BlobStore interface {
	Save(ctx context.Context, key string, r io.Reader, contentType string) (int64, error)
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
