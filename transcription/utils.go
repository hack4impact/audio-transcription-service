// Package transcription implements functions for the manipulation and
// transcription of audio files.
package transcription

import (
	"io"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/hack4impact/transcribe4all/config"
)

// SendEmail connects to an email server at host:port, switches to TLS,
// authenticates on TLS connections using the username and password, and sends
// an email from address from, to address to, with subject line subject with
// message body.
func SendEmail(username string, password string, host string, port int, to []string, subject string, body string) error {
	from := username
	auth := smtp.PlainAuth("", username, password, host)

	// The msg parameter should be an RFC 822-style email with headers first,
	// a blank line, and then the message body. The lines of msg should be CRLF
	// terminated.
	msg := []byte(msgHeaders(from, to, subject) + "\r\n" + body + "\r\n")
	addr := host + ":" + strconv.Itoa(port)
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

// ConvertAudioIntoFormat converts encoded audio into the required format.
func ConvertAudioIntoFormat(filePath, fileExt string) (string, error) {
	// http://cmusphinx.sourceforge.net/wiki/faq
	// -ar 16000 sets frequency to required 16khz
	// -ac 1 sets the number of audio channels to 1
	newPath := filePath + "." + fileExt
	cmd := exec.Command("ffmpeg", "-i", filePath, "-ar", "16000", "-ac", "1", newPath)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return newPath, nil
}

// DownloadFileFromURL locally downloads an audio file stored at url.
func DownloadFileFromURL(url string) (string, error) {
	// Taken from https://github.com/thbar/golang-playground/blob/master/download-files.go
	filePath := filePathFromURL(url)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get file contents
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Write the body to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func filePathFromURL(url string) string {
	tokens := strings.Split(url, "/")
	filePath := tokens[len(tokens)-1]
	return filePath
}

// MakeIBMTaskFunction returns a task function for transcription using IBM transcription functions.
func MakeIBMTaskFunction(audioURL string, emailAddresses []string, searchWords []string) func(string) error {
	return func(id string) error {
		filePath, err := DownloadFileFromURL(audioURL)
		if err != nil {
			return err
		}
		defer os.Remove(filePath)

		log.WithField("task", id).
			Debugf("Downloaded file at %s to %s", audioURL, filePath)

		flacPath, err := ConvertAudioIntoFormat(filePath, "flac")
		if err != nil {
			return err
		}
		defer os.Remove(flacPath)

		log.WithField("task", id).
			Debugf("Converted file %s to %s", filePath, flacPath)

		ibmResult, err := TranscribeWithIBM(flacPath, config.Config.IBMUsername, config.Config.IBMPassword)
		if err != nil {
			return err
		}
		transcript := GetTranscript(ibmResult)

		log.WithField("task", id).
			Info(transcript)

		// TODO: save data to MongoDB and file to Backblaze.

		if err := SendEmail(config.Config.EmailUsername, config.Config.EmailPassword, "smtp.gmail.com", 25, emailAddresses, "IBM Transcription Done!", transcript); err != nil {
			return err
		}

		log.WithField("task", id).
			Debugf("Sent email to %v", emailAddresses)

		return nil
	}
}
