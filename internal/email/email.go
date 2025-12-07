package email

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendHTML sends a multipart/alternative email with both text and HTML parts.
func SendHTML(to, subject, textBody, htmlBody string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, pass, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	boundary := "==BOUNDARY=="
	header := fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", to, subject, boundary)

	// Build multipart message
	msg := header
	// plain text part
	msg += fmt.Sprintf("--%s\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s\r\n\r\n", boundary, textBody)
	// html part
	msg += fmt.Sprintf("--%s\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n%s\r\n\r\n", boundary, htmlBody)
	msg += fmt.Sprintf("--%s--\r\n", boundary)

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
