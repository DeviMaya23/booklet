package storage

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/devi/booklet/internal/platform/observability"
	"github.com/stretchr/testify/require"
)

func newTestR2Storage(now *time.Time) *r2Storage {
	client := s3.New(s3.Options{
		Region:       "auto",
		BaseEndpoint: aws.String("https://acct.r2.cloudflarestorage.com"),
		Credentials:  credentials.NewStaticCredentialsProvider("key", "secret", ""),
	})
	tel := observability.NewTelemetry(nil, nil, nil)
	hist, _ := tel.Meter.Float64Histogram("test")
	return &r2Storage{
		client:             client,
		presign:            s3.NewPresignClient(client),
		bucket:             "bucket",
		tel:                tel,
		presignURLDuration: hist,
		now:                func() time.Time { return *now },
	}
}

func TestGenerateDeterministicPresignedGetURL_StableWithinWindow(t *testing.T) {
	now := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	r := newTestR2Storage(&now)

	first, err := r.GenerateDeterministicPresignedGetURL(context.Background(), "users/u/thumbnails/f.jpg")
	require.NoError(t, err)

	now = now.Add(10 * time.Hour)
	second, err := r.GenerateDeterministicPresignedGetURL(context.Background(), "users/u/thumbnails/f.jpg")
	require.NoError(t, err)

	require.Equal(t, first, second)
}

func TestGenerateDeterministicPresignedGetURL_ChangesAfterWindowRollover(t *testing.T) {
	now := time.Date(2026, 1, 1, 23, 59, 0, 0, time.UTC)
	r := newTestR2Storage(&now)

	before, err := r.GenerateDeterministicPresignedGetURL(context.Background(), "users/u/thumbnails/f.jpg")
	require.NoError(t, err)

	now = now.Add(2 * time.Minute)
	after, err := r.GenerateDeterministicPresignedGetURL(context.Background(), "users/u/thumbnails/f.jpg")
	require.NoError(t, err)

	require.NotEqual(t, before, after)
}
