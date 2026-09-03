package handler

import (
	"context"
	"time"
)

type Presigner interface {
	GeneratePresignedGetURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}
