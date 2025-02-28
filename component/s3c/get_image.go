package s3c

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	s32 "github.com/aws/aws-sdk-go/service/s3"
)

func (s *S3Component) GetImageWithExpireLink(ctx context.Context, imageKey string, duration time.Duration) (string, error) {
	req, _ := s.svc.GetObjectRequest(&s32.GetObjectInput{
		Bucket: aws.String(s.cfg.bucket),
		Key:    aws.String(imageKey),
	})

	return req.Presign(duration)
}
