package s3c

import (
	"context"
	"os"
	"path/filepath"

	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

func (s *S3Component) Download(ctx context.Context, key string) (*os.File, error) {
	s.logger.Infof("Downloading file from S3: bucket=%s, key=%s", s.cfg.bucket, key)

	fileName := filepath.Base(key)
	tmpFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp file for download")
	}

	downloader := s3manager.NewDownloaderWithClient(s.svc)

	_, err = downloader.DownloadWithContext(ctx, tmpFile, &s3.GetObjectInput{
		Bucket: &s.cfg.bucket,
		Key:    &key,
	})
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, errors.Wrapf(err, "failed to download file %q from S3", key)
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, errors.Wrap(err, "failed to seek downloaded file")
	}

	s.logger.Infof("Downloaded file: %s", tmpFile.Name())
	return tmpFile, nil
}
