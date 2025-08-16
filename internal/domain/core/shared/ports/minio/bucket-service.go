package minio

import "context"

type BucketService interface {
	CreateBucket(ctx context.Context, name string) error
	BucketExists(ctx context.Context, name string) (bool, error)
}
