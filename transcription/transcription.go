package transcription

import (
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

// SendEmail connects to an email server at host:port, switches to TLS,
// authenticates on TLS connections using the username and password, and sends
// an email from address from, to address to, with subject line subject with message body.
func SendEmail(username string, password string, host string, port int, to []string, subject string, body string) error {
	from := username
	auth := smtp.PlainAuth("", username, password, host)

	// The msg parameter should be an RFC 822-style email with headers first,
	// a blank line, and then the message body. The lines of msg should be CRLF terminated.
	msg := []byte(msgHeaders(from, to, subject) + "\r\n" + body + "\r\n")
	addr := host + ":" + string(port)
	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		return err
	}
	return nil
}

func msgHeaders(from string, to []string, subject string) string {
	fromHeader := "From: " + from
	toHeader := "To: " + strings.Join(to, ", ")
	subjectHeader := "Subject: " + subject
	msgHeaders := []string{fromHeader, toHeader, subjectHeader}
	return strings.Join(msgHeaders, "\r\n")
}

// DownloadFileFromURL locally downloads an audio file stored at url.
func DownloadFileFromURL(url string) error {
	// Taken from https://github.com/thbar/golang-playground/blob/master/download-files.go
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	// Create the file
	output, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer output.Close()

	// Get file contents
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Write the body to file
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	return nil
}
