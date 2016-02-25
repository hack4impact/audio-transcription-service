package transcription

import (
	"log"
	"net/smtp"
	"os"
	"strings"
)

func sendEmail(from string, to []string, subject string, body string) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		os.Getenv("MAIL_EMAIL"),
		os.Getenv("MAIL_PASSWORD"),
		"smtp.gmail.com",
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	msg := []byte("To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:25", auth, from, to, msg)

	if err != nil {
		log.Fatal(err)
	}
}
