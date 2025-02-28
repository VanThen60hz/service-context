package s3c

import (
	"flag"

	sctx "github.com/VanThen60hz/service-context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type S3Component struct {
	id     string
	logger sctx.Logger
	cfg    s3Config
	svc    *s3.S3
}

type s3Config struct {
	apiKey    string
	apiSecret string
	region    string
	bucket    string
	domain    string // Added domain field
}

func NewS3Component(id string) *S3Component {
	return &S3Component{id: id}
}

func (s *S3Component) ID() string {
	return s.id
}

func (s *S3Component) InitFlags() {
	flag.StringVar(&s.cfg.apiKey, "s3-api-key", "", "AWS S3 API key")
	flag.StringVar(&s.cfg.apiSecret, "s3-api-secret", "", "AWS S3 API secret")
	flag.StringVar(&s.cfg.region, "s3-region", "", "AWS S3 region")
	flag.StringVar(&s.cfg.bucket, "s3-bucket", "", "AWS S3 bucket")
	flag.StringVar(&s.cfg.domain, "s3-domain", "s3.amazonaws.com", "AWS S3 domain") // Added domain flag
}

func (s *S3Component) Activate(ctx sctx.ServiceContext) error {
	s.logger = ctx.Logger(s.id)

	if err := s.cfg.validate(); err != nil {
		return err
	}

	creds := credentials.NewStaticCredentials(s.cfg.apiKey, s.cfg.apiSecret, "")
	awsConfig := aws.NewConfig().WithRegion(s.cfg.region).WithCredentials(creds)
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS session")
	}

	s.svc = s3.New(sess)
	return nil
}

func (s *S3Component) Stop() error {
	return nil
}

func (cfg *s3Config) validate() error {
	if cfg.apiKey == "" {
		return errors.New("AWS S3 API key is missing")
	}
	if cfg.apiSecret == "" {
		return errors.New("AWS S3 API secret is missing")
	}
	if cfg.region == "" {
		return errors.New("AWS S3 region is missing")
	}
	if cfg.bucket == "" {
		return errors.New("AWS S3 bucket is missing")
	}
	return nil
}
