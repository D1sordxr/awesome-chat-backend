package minio

import (
	cfg "awesome-chat/internal/infrastructure/config/minio"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

type Connection struct {
	*minio.Client
}

func NewConnection(cfg cfg.Config) *Connection {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	return &Connection{client}
}

func (c *Connection) Start(ctx context.Context) error {
	_, err := c.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("MinIO health check failed: %w", err)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if _, err = c.ListBuckets(ctx); err != nil {
				return fmt.Errorf("MinIO health check failed: %w", err)
			}
		}
	}
}

func (c *Connection) Shutdown(_ context.Context) error {
	return nil
}
