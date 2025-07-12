package email

import (
	"log"
	"net/smtp"
	"os"
)

// SendFunc is used by Send to dispatch email messages. It can be replaced
// in tests to capture outgoing emails.
var SendFunc = func(to, subject, body string) error {
	from := smtpFrom()
	addr := smtpAddr()
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")
	if err := smtp.SendMail(addr, nil, from, []string{to}, msg); err != nil {
		return err
	}
	log.Printf("sent email to %s via %s", to, addr)
	return nil
}

func smtpAddr() string {
	if v := os.Getenv("SMTP_ADDR"); v != "" {
		return v
	}
	return "localhost:25"
}

func smtpFrom() string {
	if v := os.Getenv("SMTP_FROM"); v != "" {
		return v
	}
	return "codex@example.com"
}

// Send delivers an email using the configured SendFunc.
func Send(to, subject, body string) error {
	return SendFunc(to, subject, body)
}
