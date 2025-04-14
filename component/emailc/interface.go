package emailc

import (
	"context"
)

type Email interface {
	// SendGenericOTP sends an OTP email
	SendGenericOTP(ctx context.Context, toEmail, subject string, data OTPMailData) error
	// SendGenericLink sends a link email
	SendGenericLink(ctx context.Context, toEmail, subject string, data LinkMailData) error
}

// OTPMailData represents the data needed for OTP email template
type OTPMailData struct {
	Title         string
	UserEmail     string
	MessageIntro  string
	OTP           string
	OTPTypeDesc   string
	ExpireMinutes int
}

// LinkMailData represents the data needed for link email template
type LinkMailData struct {
	Title         string
	UserEmail     string
	MessageIntro  string
	Link          string
	ButtonText    string
	LinkTypeDesc  string
	ExpireMinutes *int
}

// Ensure EmailComponent implements Email interface
var _ Email = (*EmailComponent)(nil)
