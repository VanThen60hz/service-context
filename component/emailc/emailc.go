package emailc

import (
	"flag"

	sctx "github.com/VanThen60hz/service-context"
	"github.com/pkg/errors"
)

type EmailComponent struct {
	id     string
	logger sctx.Logger
	cfg    emailConfig
}

type emailConfig struct {
	smtpHost  string
	smtpPort  int
	emailUser string
	emailPass string
}

func NewEmailComponent(id string) *EmailComponent {
	return &EmailComponent{id: id}
}

func (e *EmailComponent) ID() string {
	return e.id
}

func (e *EmailComponent) InitFlags() {
	flag.StringVar(&e.cfg.smtpHost, "email-smtp-host", "smtp.gmail.com", "SMTP server host")
	flag.IntVar(&e.cfg.smtpPort, "email-smtp-port", 587, "SMTP server port")
	flag.StringVar(&e.cfg.emailUser, "email-user", "", "Email username")
	flag.StringVar(&e.cfg.emailPass, "email-password", "", "Email password")
}

func (e *EmailComponent) Activate(ctx sctx.ServiceContext) error {
	e.logger = ctx.Logger(e.id)

	if err := e.cfg.validate(); err != nil {
		return err
	}

	return nil
}

func (e *EmailComponent) Stop() error {
	return nil
}

func (cfg *emailConfig) validate() error {
	if cfg.emailUser == "" {
		return errors.New("email username is missing")
	}
	if cfg.emailPass == "" {
		return errors.New("email password is missing")
	}
	return nil
}
