package store

import "context"

type URLStore interface {
	Set(ctx context.Context, short string, long string) error
	Get(ctx context.Context, short string) (string, error)
	All(ctx context.Context) (map[string]string, error)
	Close() error
	Delete(ctx context.Context, short string) error
}
