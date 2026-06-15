package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
	bucket string
}

type S3Config struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

func NewS3(cfg S3Config) (*S3, error) {
	if strings.TrimSpace(cfg.Endpoint) == "" {
		return nil, fmt.Errorf("S3 endpoint is required")
	}
	if strings.TrimSpace(cfg.Bucket) == "" {
		return nil, fmt.Errorf("S3 bucket is required")
	}
	if strings.TrimSpace(cfg.AccessKey) == "" || strings.TrimSpace(cfg.SecretKey) == "" {
		return nil, fmt.Errorf("S3 credentials are required")
	}
	region := strings.TrimSpace(cfg.Region)
	if region == "" {
		region = "auto"
	}

	client := s3.New(s3.Options{
		Region: region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
		BaseEndpoint: aws.String(strings.TrimRight(cfg.Endpoint, "/")),
		UsePathStyle: true,
	})

	return &S3{client: client, bucket: cfg.Bucket}, nil
}

func (s *S3) Save(ctx context.Context, key string, r io.Reader, contentType string) (int64, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return 0, fmt.Errorf("read body: %w", err)
	}
	in := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	}
	if contentType != "" {
		in.ContentType = aws.String(contentType)
	}
	if _, err := s.client.PutObject(ctx, in); err != nil {
		return 0, fmt.Errorf("put object: %w", err)
	}
	return int64(len(body)), nil
}

func (s *S3) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	return out.Body, nil
}

func (s *S3) Delete(ctx context.Context, key string) error {
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

func S3ConfigFromEnv() S3Config {
	return S3Config{
		Endpoint:  firstEnv("AWS_ENDPOINT_URL", "S3_ENDPOINT", "ENDPOINT"),
		Region:    firstEnv("AWS_DEFAULT_REGION", "S3_REGION", "REGION"),
		Bucket:    firstEnv("AWS_S3_BUCKET_NAME", "S3_BUCKET", "BUCKET"),
		AccessKey: firstEnv("AWS_ACCESS_KEY_ID", "S3_ACCESS_KEY_ID", "ACCESS_KEY_ID"),
		SecretKey: firstEnv("AWS_SECRET_ACCESS_KEY", "S3_SECRET_ACCESS_KEY", "SECRET_ACCESS_KEY"),
	}
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}
