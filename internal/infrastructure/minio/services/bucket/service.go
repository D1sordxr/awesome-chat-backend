package bucket

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"

	appPorts "awesome-chat/internal/domain/app/ports"
	root "awesome-chat/internal/infrastructure/minio"
)

type Service struct {
	log appPorts.Logger
	*root.Connection
}

func NewService(
	log appPorts.Logger,
	conn *root.Connection,
) *Service {
	return &Service{log: log, Connection: conn}
}

func (s *Service) CreateBucket(ctx context.Context, name string) error {
	const op = "minio.BucketService.CreateBucket"

	ok, err := s.Client.BucketExists(ctx, name)
	switch {
	case err != nil:
		return fmt.Errorf("%s: %w", op, err)
	case ok:
		s.log.Debug(fmt.Sprintf("%s: bucket '%s' already exists", op, name))
		return nil
	default:
		err = s.Client.MakeBucket(ctx, name, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("%s: failed to create bucket: %w", op, err)
		}
		s.log.Info(fmt.Sprintf("%s: bucket '%s' created successfully", op, name))
	}

	return nil
}

func (s *Service) BucketExists(ctx context.Context, name string) (bool, error) {
	return s.Client.BucketExists(ctx, name)
}
