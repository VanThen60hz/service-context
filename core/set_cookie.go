package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookieName = "accessToken"
	defaultCookieMaxAge   = 604800 // 7 days in seconds
	defaultCookiePath     = "/"
)

// CookieOptions defines the configuration options for cookies
type CookieOptions struct {
	MaxAge   int
	Domain   string
	Path     string
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

// DefaultCookieOptions returns the default cookie configuration
func DefaultCookieOptions(path string) CookieOptions {
	if path == "" {
		path = defaultCookiePath
	}
	return CookieOptions{
		MaxAge:   defaultCookieMaxAge,
		Domain:   "",
		Path:     path,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}

// mergeCookieOptions merges the default options with custom options
func mergeCookieOptions(defaultOpts CookieOptions, customOpts ...CookieOptions) CookieOptions {
	if len(customOpts) == 0 {
		return defaultOpts
	}

	opts := defaultOpts
	custom := customOpts[0]

	if custom.Path != "" {
		opts.Path = custom.Path
	}
	if custom.MaxAge != 0 {
		opts.MaxAge = custom.MaxAge
	}
	if custom.Domain != "" {
		opts.Domain = custom.Domain
	}
	if custom.Secure {
		opts.Secure = custom.Secure
	}
	if custom.HttpOnly {
		opts.HttpOnly = custom.HttpOnly
	}
	if custom.SameSite != 0 {
		opts.SameSite = custom.SameSite
	}

	return opts
}

// SetCookie sets a cookie with the given name, value and options (production)
func SetCookie(c *gin.Context, name, value string, opts CookieOptions) {
	secure := opts.Secure
	if secure {
		secure = IsHTTPS(c)
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     opts.Path,
		Domain:   opts.Domain,
		MaxAge:   opts.MaxAge,
		HttpOnly: opts.HttpOnly,
		Secure:   secure,
		SameSite: opts.SameSite,
	})
}

// SetDevelopmentCookie sets a cookie for development (localhost, HTTP)
func SetDevelopmentCookie(c *gin.Context, name, value string, opts CookieOptions) {
	opts.Secure = false
	if opts.SameSite == http.SameSiteNoneMode {
		opts.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     opts.Path,
		Domain:   opts.Domain,
		MaxAge:   opts.MaxAge,
		HttpOnly: opts.HttpOnly,
		Secure:   opts.Secure,
		SameSite: opts.SameSite,
	})
}

// SetAccessTokenCookie sets the access token cookie with standard configuration
func SetAccessTokenCookie(c *gin.Context, token string, path string, opts ...CookieOptions) {
	options := mergeCookieOptions(DefaultCookieOptions(path), opts...)
	SetCookie(c, accessTokenCookieName, token, options)
}

// SetAccessTokenCookieWithDefaultPath sets the access token cookie with default path
func SetAccessTokenCookieWithDefaultPath(c *gin.Context, token string, opts ...CookieOptions) {
	SetAccessTokenCookie(c, token, defaultCookiePath, opts...)
}

// SetAccessTokenCookieForDevelopment sets the access token cookie for local development
func SetAccessTokenCookieForDevelopment(c *gin.Context, token string, path string, opts ...CookieOptions) {
	options := mergeCookieOptions(DefaultCookieOptions(path), opts...)
	options.Secure = false
	if options.SameSite == http.SameSiteNoneMode {
		options.SameSite = http.SameSiteLaxMode
	}
	SetDevelopmentCookie(c, accessTokenCookieName, token, options)
}

// ClearAccessTokenCookie clears the access token cookie
func ClearAccessTokenCookie(c *gin.Context, path string, opts ...CookieOptions) {
	options := mergeCookieOptions(DefaultCookieOptions(path), opts...)
	options.MaxAge = -1
	SetCookie(c, accessTokenCookieName, "", options)
}

// ClearAccessTokenCookieWithDefaultPath clears the access token cookie with default path
func ClearAccessTokenCookieWithDefaultPath(c *gin.Context, opts ...CookieOptions) {
	ClearAccessTokenCookie(c, defaultCookiePath, opts...)
}
