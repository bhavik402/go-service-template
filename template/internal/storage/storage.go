// Package storage is an object storage abstraction (S3-compatible;
// MinIO for local development via docker-compose).
package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
