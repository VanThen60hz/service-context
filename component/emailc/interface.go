package emailc

import (
	"context"
)

type Email interface {
	// SendGenericOTP sends an OTP email
	SendGenericOTP(ctx context.Context, toEmail, subject string, data OTPMailData) error
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

// Ensure EmailComponent implements Email interface
var _ Email = (*EmailComponent)(nil)
