package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/riverqueue/river"
	rivertype "github.com/riverqueue/river/rivertype"
)

var ErrCharacterNotOwned = errors.New("one or more character IDs do not belong to the user")

const (
	PresignTTL    = 15 * time.Minute
	PresignGetTTL = 1 * time.Hour
)

type StorageService interface {
	GeneratePresignedPutURL(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

type JobInserter interface {
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}
