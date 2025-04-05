package s3c

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

func (s *S3Component) GetPresignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	req, _ := s.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(duration)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate presigned URL")
	}

	return url, nil
}
