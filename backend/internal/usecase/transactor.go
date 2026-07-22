package usecase

import "context"

type Transactor interface {
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
