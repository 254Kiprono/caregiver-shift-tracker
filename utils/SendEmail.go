package utils

import (
	"fmt"
	"net/smtp"
)

// SendEmail sends an email using the configured SMTP server
func SendEmail(email, subject, body string) error {
	smtpServer := cfg.SMTPServer
	smtpPort := cfg.SMTPPort
	smtpUsername := cfg.SenderEmail
	smtpPassword := cfg.EmailPassword

	if smtpPort == "" {
		smtpPort = "587"
	}

	// Email headers
	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		smtpUsername, email, subject, body,
	)

	// Authentication for SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	// Send the email
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUsername, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
