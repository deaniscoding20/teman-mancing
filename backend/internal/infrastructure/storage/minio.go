package storage

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"

	"teman-mancing/backend/internal/config"
)

// MinIOClient wraps the MinIO S3 SDK client.
type MinIOClient struct {
	Client  *minio.Client
	Buckets config.MinIOBuckets
	logger  *zap.Logger
}

// NewMinIOClient connects to the MinIO instance and ensures required buckets exist.
func NewMinIOClient(cfg *config.Config, log *zap.Logger) (*MinIOClient, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	mc := &MinIOClient{
		Client:  client,
		Buckets: cfg.MinIOBuckets,
		logger:  log,
	}

	return mc, nil
}

// Ping tests MinIO connectivity.
func (m *MinIOClient) Ping(ctx context.Context) error {
	if m == nil || m.Client == nil {
		return fmt.Errorf("minio client is nil")
	}
	_, err := m.Client.ListBuckets(ctx)
	return err
}

// EnsureBuckets checks for bucket existence and creates them if absent.
func (m *MinIOClient) EnsureBuckets(ctx context.Context) error {
	bucketList := []string{
		m.Buckets.Avatars,
		m.Buckets.Spots,
		m.Buckets.Events,
	}

	for _, bucket := range bucketList {
		if strings.TrimSpace(bucket) == "" {
			continue
		}
		exists, err := m.Client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("failed to check bucket '%s': %w", bucket, err)
		}
		if !exists {
			err = m.Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return fmt.Errorf("failed to make bucket '%s': %w", bucket, err)
			}
			m.logger.Info("MinIO bucket created", zap.String("bucket", bucket))
		}
	}
	return nil
}

// PresignedGetObject generates a temporary pre-signed URL to view a private file.
func (m *MinIOClient) PresignedGetObject(ctx context.Context, bucket, objectKey string, expires time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	return m.Client.PresignedGetObject(ctx, bucket, objectKey, expires, reqParams)
}
