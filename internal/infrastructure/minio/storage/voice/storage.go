package voice

import (
	minioPorts "awesome-chat/internal/domain/core/shared/ports/minio"
	root "awesome-chat/internal/infrastructure/minio"
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"sync"
)

type Storage struct {
	conn       *root.Connection
	bucketName string
	bucketSvc  minioPorts.BucketService
	once       sync.Once
}

func NewStorage(
	conn *root.Connection,
	bucketName string,
	bucketSvc minioPorts.BucketService,
) *Storage {
	return &Storage{
		conn:       conn,
		bucketName: bucketName,
		bucketSvc:  bucketSvc,
	}
}

func (s *Storage) Add(
	ctx context.Context,
	objID string,
	data []byte,
) (err error) {
	s.once.Do(func() {
		exists, bucketErr := s.bucketSvc.BucketExists(ctx, s.bucketName)
		if bucketErr != nil {
			err = fmt.Errorf("failed to check bucket existence: %w", bucketErr)
			return
		}
		if !exists {
			if bucketErr = s.bucketSvc.CreateBucket(ctx, s.bucketName); bucketErr != nil {
				err = fmt.Errorf("failed to create bucket: %w", bucketErr)
			}
		}
	})
	if err != nil {
		return err
	}

	if _, err = s.conn.PutObject(
		ctx,
		s.bucketName,
		objID,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{},
	); err != nil {
		return fmt.Errorf("failed to upload voice message: %w", err)
	}

	return nil
}
