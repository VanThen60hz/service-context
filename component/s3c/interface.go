package s3c

import (
	"context"
	"os"
	"time"
)

type S3 interface {
	UploadFileData(ctx context.Context, fileData []byte, cloudFolder string) (string, error)
	// Upload image to AWS S3 and response URL
	Upload(ctx context.Context, fileName string, cloudFolder string) (string, error)
	// Get image link from uploaded with imageKey and duration
	GetImageWithExpireLink(ctx context.Context, imageKey string, duration time.Duration) (string, error)
	// Delete image with imageKey and duration
	DeleteImages(ctx context.Context, fileKeys []string) error
	// Delete any object
	DeleteObject(ctx context.Context, key string) error
	// Get presigned URL for downloading file
	GetPresignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
	// Download file from S3
	Download(ctx context.Context, key string) (*os.File, error)
}
