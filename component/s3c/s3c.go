package s3c

import (
	"flag"
	"fmt"

	sctx "github.com/VanThen60hz/service-context"

	"github.com/VanThen60hz/service-context/component/logger"
	"github.com/VanThen60hz/service-context/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s32 "github.com/aws/aws-sdk-go/service/s3"
)

type s3 struct {
	id     string
	prefix string
	logger logger.Logger

	cfg s3Config

	session *session.Session
	service *s32.S3
}

type s3Config struct {
	s3ApiKey    string
	s3ApiSecret string
	s3Region    string
	s3Bucket    string
}

func NewS3(id string, prefix ...string) *s3 {
	pre := "aws-s3"
	if len(prefix) > 0 {
		pre = prefix[0]
	}

	return &s3{
		id:     id,
		prefix: pre,
	}
}

func (s *s3) ID() string {
	return s.id
}

func (s *s3) InitFlags() {
	flag.StringVar(&s.cfg.s3ApiKey, fmt.Sprintf("%s-%s", s.GetPrefix(), "api-key"), "", "S3 API key")
	flag.StringVar(&s.cfg.s3ApiSecret, fmt.Sprintf("%s-%s", s.GetPrefix(), "api-secret"), "", "S3 API secret key")
	flag.StringVar(&s.cfg.s3Region, fmt.Sprintf("%s-%s", s.GetPrefix(), "region"), "", "S3 region")
	flag.StringVar(&s.cfg.s3Bucket, fmt.Sprintf("%s-%s", s.GetPrefix(), "bucket"), "", "S3 bucket")
}

func (s *s3) Activate(_ sctx.ServiceContext) error {
	s.logger = logger.GetCurrent().GetLogger(s.ID())

	if err := s.cfg.check(); err != nil {
		s.logger.Errorln(err)
		return err
	}

	credential := credentials.NewStaticCredentials(s.cfg.s3ApiKey, s.cfg.s3ApiSecret, "")
	_, err := credential.Get()
	if err != nil {
		s.logger.Errorln(err)
		return err
	}

	config := aws.NewConfig().WithRegion(s.cfg.s3Region).WithCredentials(credential)
	ss, err := session.NewSession(config)
	if err != nil {
		s.logger.Errorln(err)
		return err
	}

	s.service = s32.New(ss, config)
	s.session = ss

	return nil
}

func (s *s3) Stop() error {
	// Implement any cleanup logic if necessary
	return nil
}

func (s *s3) GetPrefix() string {
	return s.prefix
}

func (cfg *s3Config) check() error {
	if len(cfg.s3ApiKey) < 1 {
		return core.ErrS3ApiKeyMissing
	}
	if len(cfg.s3ApiSecret) < 1 {
		return core.ErrS3ApiSecretKeyMissing
	}
	if len(cfg.s3Bucket) < 1 {
		return core.ErrS3BucketMissing
	}
	if len(cfg.s3Region) < 1 {
		return core.ErrS3RegionMissing
	}
	return nil
}
