package s3

import "context"

type URLService interface {
	GenerateURL(ctx context.Context, objID string) (string, error)
	DownloadFromURL(ctx context.Context, url string) ([]byte, error)
	Read(ctx context.Context, objID string) ([]byte, error)
}
