package handler

import (
	"context"
	"time"
)

type Presigner interface {
	GeneratePresignedGetURL(ctx context.Context, key string, ttl time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key, filename string, ttl time.Duration) (string, error)
}
