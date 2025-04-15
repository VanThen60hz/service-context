package oauthc

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"

	sctx "github.com/VanThen60hz/service-context"
	"github.com/pkg/errors"
)

type OAuthComponent struct {
	id     string
	logger sctx.Logger
	cfg    oauthConfig
	state  string
}

type oauthConfig struct {
	googleClientID     string
	googleClientSecret string
	googleRedirectURL  string
	googleUserInfoURL  string

	facebookClientID     string
	facebookClientSecret string
	facebookRedirectURL  string
	facebookGraphMeURL   string
}

func NewOAuthComponent(id string) *OAuthComponent {
	return &OAuthComponent{id: id}
}

func (o *OAuthComponent) ID() string {
	return o.id
}

func (o *OAuthComponent) InitFlags() {
	flag.StringVar(&o.cfg.googleClientID, "google-client-id", "", "Google OAuth client ID")
	flag.StringVar(&o.cfg.googleClientSecret, "google-client-secret", "", "Google OAuth client secret")
	flag.StringVar(&o.cfg.googleRedirectURL, "google-redirect-url", "", "Google OAuth redirect URL")
	flag.StringVar(&o.cfg.googleUserInfoURL, "google-user-info-url", "https://www.googleapis.com/oauth2/v2/userinfo", "Google user info URL")

	flag.StringVar(&o.cfg.facebookClientID, "facebook-client-id", "", "Facebook OAuth client ID")
	flag.StringVar(&o.cfg.facebookClientSecret, "facebook-client-secret", "", "Facebook OAuth client secret")
	flag.StringVar(&o.cfg.facebookRedirectURL, "facebook-redirect-url", "", "Facebook OAuth redirect URL")
	flag.StringVar(&o.cfg.facebookGraphMeURL, "facebook-graph-me-url", "https://graph.facebook.com/me", "Facebook Graph API me URL")
}

func (o *OAuthComponent) Activate(ctx sctx.ServiceContext) error {
	o.logger = ctx.Logger(o.id)

	if err := o.cfg.validate(); err != nil {
		return err
	}

	// Generate OAuth state string
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		o.state = fmt.Sprintf("%s-fallback-state-string", o.id)
		o.logger.Error("Error generating random state string", err)
	} else {
		o.state = base64.URLEncoding.EncodeToString(b)
		o.logger.Info("Generated OAuth state string", o.state)
	}

	return nil
}

func (o *OAuthComponent) Stop() error {
	return nil
}

func (cfg *oauthConfig) validate() error {
	if cfg.googleClientID == "" {
		return errors.New("Google client ID is missing")
	}
	if cfg.googleClientSecret == "" {
		return errors.New("Google client secret is missing")
	}
	if cfg.googleRedirectURL == "" {
		return errors.New("Google redirect URL is missing")
	}

	if cfg.facebookClientID == "" {
		return errors.New("Facebook client ID is missing")
	}
	if cfg.facebookClientSecret == "" {
		return errors.New("Facebook client secret is missing")
	}
	if cfg.facebookRedirectURL == "" {
		return errors.New("Facebook redirect URL is missing")
	}

	return nil
}
