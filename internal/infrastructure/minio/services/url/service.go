package url

import (
	cachePort "awesome-chat/internal/domain/core/shared/ports/cache"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
)

type Service struct {
	minioClient *minio.Client
	bucketName  string
	urlCache    cachePort.Storage
	ttl         time.Duration
}

const defaultTTL = 15 * time.Minute

func NewUrlService(
	minioClient *minio.Client,
	bucketName string,
	urlCache cachePort.Storage,
) *Service {
	return &Service{
		minioClient: minioClient,
		bucketName:  bucketName,
		urlCache:    urlCache,
		ttl:         defaultTTL,
	}
}

// Read returns file data by objID, using cached URL or generating a new one.
func (s *Service) Read(ctx context.Context, objID string) ([]byte, error) {
	// 1. Try to get URL from cache
	cachedURL, err := s.urlCache.Read(ctx, objID)
	if err == nil {
		data, dErr := s.downloadFromURL(ctx, cachedURL)
		if dErr == nil {
			return data, nil
		}
		// If URL didn't work - delete from cache and try to generate new one
		_ = s.urlCache.Delete(ctx, objID)
	}

	// 2. Generate new pre-signed URL
	newURL, err := s.minioClient.PresignedGetObject(
		ctx,
		s.bucketName,
		objID,
		s.ttl,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	// 3. Save URL to cache
	urlStr := newURL.String()
	if err := s.urlCache.Set(ctx, objID, urlStr, s.ttl); err != nil {
		return nil, fmt.Errorf("failed to cache URL: %w", err)
	}

	// 4. Download file
	data, err := s.downloadFromURL(ctx, urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return data, nil
}

// downloadFromURL downloads file by URL (could be HTTP client with retry).
func (s *Service) downloadFromURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch file: " + resp.Status)
	}

	return io.ReadAll(resp.Body)
}
