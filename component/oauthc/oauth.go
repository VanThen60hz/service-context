package oauthc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/VanThen60hz/service-context/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

func (o *OAuthComponent) GetGoogleAuthURL(state string) string {
	config := &oauth2.Config{
		ClientID:     o.cfg.googleClientID,
		ClientSecret: o.cfg.googleClientSecret,
		RedirectURL:  o.cfg.googleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return config.AuthCodeURL(state)
}

func (o *OAuthComponent) GetFacebookAuthURL(state string) string {
	config := &oauth2.Config{
		ClientID:     o.cfg.facebookClientID,
		ClientSecret: o.cfg.facebookClientSecret,
		RedirectURL:  o.cfg.facebookRedirectURL,
		Scopes: []string{
			"public_profile",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("auth_type", "rerequest"),
	}
	return config.AuthCodeURL(state, opts...)
}

func (o *OAuthComponent) ProcessGoogleCallback(ctx context.Context, code string, state string) (*OAuthUserInfo, error) {
	if state != o.state {
		return nil, core.ErrInternalServerError.WithError("invalid state parameter")
	}

	config := &oauth2.Config{
		ClientID:     o.cfg.googleClientID,
		ClientSecret: o.cfg.googleClientSecret,
		RedirectURL:  o.cfg.googleRedirectURL,
		Endpoint:     google.Endpoint,
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError("code exchange failed").WithDebug(err.Error())
	}

	resp, err := http.Get(o.cfg.googleUserInfoURL + "?access_token=" + token.AccessToken)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError("failed getting user info").WithDebug(err.Error())
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, core.ErrInternalServerError.WithError("failed decoding user info").WithDebug(err.Error())
	}

	return &OAuthUserInfo{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		Provider:  "google",
	}, nil
}

func (o *OAuthComponent) ProcessFacebookCallback(ctx context.Context, code string, state string) (*OAuthUserInfo, error) {
	if state != o.state {
		return nil, core.ErrInternalServerError.WithError("invalid state parameter")
	}

	config := &oauth2.Config{
		ClientID:     o.cfg.facebookClientID,
		ClientSecret: o.cfg.facebookClientSecret,
		RedirectURL:  o.cfg.facebookRedirectURL,
		Endpoint:     facebook.Endpoint,
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError("code exchange failed").WithDebug(err.Error())
	}

	resp, err := http.Get(o.cfg.facebookGraphMeURL + "?fields=id,name,email&access_token=" + token.AccessToken)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError("failed getting user info").WithDebug(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError("failed reading response").WithDebug(err.Error())
	}

	var userInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, core.ErrInternalServerError.WithError("failed parsing user info").WithDebug(err.Error())
	}

	// Split full name into first and last name
	names := strings.Fields(userInfo.Name)
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}

	return &OAuthUserInfo{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		FirstName: firstName,
		LastName:  lastName,
		Provider:  "facebook",
	}, nil
}
