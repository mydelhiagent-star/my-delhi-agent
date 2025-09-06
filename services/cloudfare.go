// services/cloudflare_r2.go
package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// CloudflareR2Service handles Cloudflare R2 operations
type CloudflareR2Service struct {
	client     *s3.Client
	bucketName string
	accountID  string
}

// R2Config holds configuration for Cloudflare R2
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// NewCloudflareR2Service creates a new Cloudflare R2 service
func NewCloudflareR2Service(AccountID string, AccessKeyID string, AccessKeySecret string, BucketName string) (*CloudflareR2Service, error) {

	// Validate required fields
	if AccountID == "" || AccessKeyID == "" ||
		AccessKeySecret == "" || BucketName == "" {
		return nil, fmt.Errorf("missing required Cloudflare R2 configuration")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AccessKeyID, AccessKeySecret, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with R2 endpoint
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", AccountID))
	})

	return &CloudflareR2Service{
		client:     client,
		bucketName: BucketName,
		accountID:  AccountID,
	}, nil
}

// GeneratePresignedURL generates a presigned URL for uploading objects
func (s *CloudflareR2Service) GeneratePresignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	if expires <= 0 {
		expires = 1 * time.Hour // Default to 1 hour
	}

	// Create presign client
	presignClient := s3.NewPresignClient(s.client)

	// Generate presigned PUT URL for upload
	presignResult, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}


	return presignResult.URL, nil
}

// GeneratePresignedGetURL generates a presigned URL for downloading objects
func (s *CloudflareR2Service) GeneratePresignedGetURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	if expires <= 0 {
		expires = 1 * time.Hour // Default to 1 hour
	}

	// Create presign client
	presignClient := s3.NewPresignClient(s.client)

	// Generate presigned GET URL for download
	presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned GET URL: %w", err)
	}


	return presignResult.URL, nil
}

// UploadObject uploads an object directly to R2
func (s *CloudflareR2Service) UploadObject(ctx context.Context, key string, data []byte, contentType string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	// Set default content type if not provided
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	log.Printf("Successfully uploaded object: %s", key)
	return nil
}

// DeleteObject deletes an object from R2
func (s *CloudflareR2Service) DeleteObject(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	log.Printf("Successfully deleted object: %s", key)
	return nil
}

// ListObjects lists objects in the bucket
func (s *CloudflareR2Service) ListObjects(ctx context.Context, prefix string, maxKeys int32) ([]string, error) {
	if maxKeys <= 0 {
		maxKeys = 1000 // Default max keys
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		// MaxKeys: maxKeys,
	}

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var keys []string
	for _, object := range result.Contents {
		keys = append(keys, *object.Key)
	}

	return keys, nil
}

// GetObjectInfo gets information about an object
func (s *CloudflareR2Service) GetObjectInfo(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object info: %w", err)
	}

	return result, nil
}

// GetBucketInfo gets information about the bucket
func (s *CloudflareR2Service) GetBucketInfo(ctx context.Context) (*s3.HeadBucketOutput, error) {
	result, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get bucket info: %w", err)
	}

	return result, nil
}
