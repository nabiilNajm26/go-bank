package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

const (
	MaxFileSize     = 5 * 1024 * 1024 // 5MB max
	MaxFilesPerUser = 100             // Safety limit
)

type S3Service struct {
	client     *s3.Client
	uploader   *manager.Uploader
	bucketName string
}

type UploadResult struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

func NewS3Service() (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME environment variable is required")
	}

	return &S3Service{
		client:     client,
		uploader:   uploader,
		bucketName: bucketName,
	}, nil
}

func (s *S3Service) UploadProfileImage(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	// Validate file size
	if header.Size > MaxFileSize {
		return nil, fmt.Errorf("file size %d exceeds maximum allowed size %d", header.Size, MaxFileSize)
	}

	// Validate file type
	if !isValidImageType(header.Filename) {
		return nil, fmt.Errorf("invalid file type. Only PNG, JPG, JPEG allowed")
	}

	// Check user file count (safety measure)
	count, err := s.getUserFileCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user file count: %v", err)
	}
	if count >= MaxFilesPerUser {
		return nil, fmt.Errorf("user has reached maximum file limit (%d)", MaxFilesPerUser)
	}

	// Generate unique key
	ext := filepath.Ext(header.Filename)
	key := fmt.Sprintf("profile-images/%s/%s%s", userID.String(), uuid.New().String(), ext)

	// Get content type
	contentType := "application/octet-stream"
	if len(header.Header["Content-Type"]) > 0 {
		contentType = header.Header["Content-Type"][0]
	}

	// Upload to S3
	result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &key,
		Body:        file,
		ContentType: &contentType,
		Metadata: map[string]string{
			"user-id":     userID.String(),
			"upload-time": time.Now().Format(time.RFC3339),
			"filename":    header.Filename,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	return &UploadResult{
		URL:      result.Location,
		Key:      key,
		Size:     header.Size,
		MimeType: contentType,
	}, nil
}

func (s *S3Service) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	return err
}

func (s *S3Service) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (s *S3Service) getUserFileCount(ctx context.Context, userID uuid.UUID) (int, error) {
	prefix := fmt.Sprintf("profile-images/%s/", userID.String())
	
	resp, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &s.bucketName,
		Prefix: &prefix,
	})
	if err != nil {
		return 0, err
	}

	if resp.KeyCount != nil {
		return int(*resp.KeyCount), nil
	}
	return 0, nil
}

func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".png", ".jpg", ".jpeg"}
	
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// SetupBucketLifecycle sets up automatic deletion after 30 days for cost control
func (s *S3Service) SetupBucketLifecycle(ctx context.Context) error {
	ruleID := "delete-old-files"
	prefix := "profile-images/"
	days := int32(30)
	
	_, err := s.client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
		Bucket: &s.bucketName,
		LifecycleConfiguration: &types.BucketLifecycleConfiguration{
			Rules: []types.LifecycleRule{
				{
					ID:     &ruleID,
					Status: types.ExpirationStatusEnabled,
					Filter: &types.LifecycleRuleFilter{
						Prefix: &prefix,
					},
					Expiration: &types.LifecycleExpiration{
						Days: &days,
					},
				},
			},
		},
	})
	
	return err
}