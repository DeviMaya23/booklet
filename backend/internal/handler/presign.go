package handler

import (
	"context"
	"time"
)

type Presigner interface {
	GeneratePresignedGetURL(ctx context.Context, key string, ttl time.Duration) (string, error)
	GenerateDeterministicPresignedGetURL(ctx context.Context, key string) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key, filename string, ttl time.Duration) (string, error)
}
