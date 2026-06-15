package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	Dir string
}

func NewLocal(dir string) (*Local, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir media dir: %w", err)
	}
	return &Local{Dir: dir}, nil
}

func (l *Local) Save(_ context.Context, key string, r io.Reader, _ string) (int64, error) {
	path := filepath.Join(l.Dir, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return 0, fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}
	defer func() { _ = f.Close() }()
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, fmt.Errorf("write file: %w", err)
	}
	return n, nil
}

func (l *Local) Open(_ context.Context, key string) (io.ReadCloser, error) {
	path := filepath.Join(l.Dir, key)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return f, nil
}

func (l *Local) Delete(_ context.Context, key string) error {
	path := filepath.Join(l.Dir, key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}
