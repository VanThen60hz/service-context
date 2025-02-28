package s3c

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	s32 "github.com/aws/aws-sdk-go/service/s3"
)

func (s *S3Component) DeleteImages(ctx context.Context, fileKeys []string) error {
	del := &s32.Delete{
		Objects: toOIDs(fileKeys),
		Quiet:   aws.Bool(false),
	}

	doi := &s32.DeleteObjectsInput{
		Bucket: aws.String(s.cfg.bucket),
		Delete: del,
	}

	res, err := s.svc.DeleteObjects(doi)
	if err != nil {
		return err
	}

	s.logger.Infoln(res)

	return nil
}

func toOIDs(keys []string) []*s32.ObjectIdentifier {
	ret := make([]*s32.ObjectIdentifier, len(keys))
	for i := 0; i < len(ret); i++ {
		oid := &s32.ObjectIdentifier{
			Key: &(keys[i]),
		}
		ret[i] = oid
	}
	return ret
}

func (s *S3Component) DeleteObject(ctx context.Context, key string) error {
	input := &s32.DeleteObjectInput{
		Bucket: aws.String(s.cfg.bucket),
		Key:    aws.String(key),
	}

	res, err := s.svc.DeleteObject(input)
	if err != nil {
		return err
	}

	s.logger.Infoln(res)

	return nil
}
