package storage

import (
	"fmt"
	"os"
	"strings"
)

func NewFromEnv() (BlobStore, string, error) {
	driver := strings.ToLower(strings.TrimSpace(os.Getenv("MEDIA_STORAGE")))
	if driver == "" {
		if cfg := S3ConfigFromEnv(); cfg.Endpoint != "" && cfg.Bucket != "" {
			driver = "s3"
		} else {
			driver = "local"
		}
	}

	switch driver {
	case "s3", "railway":
		cfg := S3ConfigFromEnv()
		store, err := NewS3(cfg)
		if err != nil {
			return nil, "", err
		}
		return store, "s3", nil
	case "local":
		dir := strings.TrimSpace(os.Getenv("MEDIA_STORAGE_DIR"))
		if dir == "" {
			dir = "./data/media"
		}
		store, err := NewLocal(dir)
		if err != nil {
			return nil, "", err
		}
		return store, "local", nil
	default:
		return nil, "", fmt.Errorf("unknown MEDIA_STORAGE driver %q (use local or s3)", driver)
	}
}
