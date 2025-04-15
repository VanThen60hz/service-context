package oauthc

import (
	"context"

	"golang.org/x/oauth2"
)

type OAuth interface {
	// GetGoogleAuthURL returns the Google OAuth URL
	GetGoogleAuthURL(state string) string
	// GetFacebookAuthURL returns the Facebook OAuth URL
	GetFacebookAuthURL(state string) string
	// ProcessGoogleCallback processes the Google OAuth callback
	ProcessGoogleCallback(ctx context.Context, code string, state string) (*OAuthUserInfo, error)
	// ProcessFacebookCallback processes the Facebook OAuth callback
	ProcessFacebookCallback(ctx context.Context, code string, state string) (*OAuthUserInfo, error)
	// GetStateString returns the current OAuth state string
	GetStateString() string
}

// OAuthUserInfo represents the user information returned from OAuth providers
type OAuthUserInfo struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Provider  string // "google" or "facebook"
}

// OAuthConfig represents the configuration for OAuth providers
type OAuthConfig struct {
	GoogleConfig   oauth2.Config
	FacebookConfig oauth2.Config
}

// Ensure OAuthComponent implements OAuth interface
var _ OAuth = (*OAuthComponent)(nil)
