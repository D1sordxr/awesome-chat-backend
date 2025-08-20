package s3

import "context"

type Storage interface {
	Add(ctx context.Context, objID string, data []byte) (err error)
	Delete(ctx context.Context, objID string) (err error)
}
