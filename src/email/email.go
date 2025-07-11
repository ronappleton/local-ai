package email

import "log"

// SendFunc is used by Send to dispatch email messages. It can be replaced
// in tests to capture outgoing emails.
var SendFunc = func(to, subject, body string) error {
	log.Printf("sending email to %s: %s - %s", to, subject, body)
	return nil
}

// Send delivers an email using the configured SendFunc.
func Send(to, subject, body string) error {
	return SendFunc(to, subject, body)
}
