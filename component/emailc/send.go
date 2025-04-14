package emailc

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
)

func (e *EmailComponent) SendGenericOTP(ctx context.Context, toEmail, subject string, data OTPMailData) error {
	auth := smtp.PlainAuth("", e.cfg.emailUser, e.cfg.emailPass, e.cfg.smtpHost)

	// Parse embedded template
	tmpl, err := template.New("otp").Parse(otpTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Email headers
	headers := map[string]string{
		"From":             e.cfg.emailUser,
		"To":               toEmail,
		"Subject":          subject,
		"MIME-Version":     "1.0",
		"Content-Type":     "text/html; charset=\"utf-8\"",
		"List-Unsubscribe": fmt.Sprintf("<%s>", e.cfg.emailUser),
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	// Send email
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", e.cfg.smtpHost, e.cfg.smtpPort),
		auth,
		e.cfg.emailUser,
		[]string{toEmail},
		[]byte(message),
	)
}

func (e *EmailComponent) SendGenericLink(ctx context.Context, toEmail, subject string, data LinkMailData) error {
	auth := smtp.PlainAuth("", e.cfg.emailUser, e.cfg.emailPass, e.cfg.smtpHost)

	// Parse embedded template
	tmpl, err := template.New("link").Parse(linkTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Email headers
	headers := map[string]string{
		"From":             e.cfg.emailUser,
		"To":               toEmail,
		"Subject":          subject,
		"MIME-Version":     "1.0",
		"Content-Type":     "text/html; charset=\"utf-8\"",
		"List-Unsubscribe": fmt.Sprintf("<%s>", e.cfg.emailUser),
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	// Send email
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", e.cfg.smtpHost, e.cfg.smtpPort),
		auth,
		e.cfg.emailUser,
		[]string{toEmail},
		[]byte(message),
	)
}
