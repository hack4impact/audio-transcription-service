// Package transcription implements functions for the manipulation and
// transcription of audio files.
package transcription

import (
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jordan-wright/email"
	"github.com/juju/errors"

	"github.com/hack4impact/transcribe4all/config"
)

// SendEmail connects to an email server at host:port and sends an email from
// address from, to address to, with subject line subject with message body.
func SendEmail(username string, password string, host string, port int, to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", username, password, host)
	addr := host + ":" + strconv.Itoa(port)

	message := email.Email{
		From:    username,
		To:      to,
		Subject: subject,
		Text:    []byte(body),
	}
	if err := message.Send(addr, auth); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// ConvertAudioIntoFormat converts encoded audio into the required format.
func ConvertAudioIntoFormat(filePath, fileExt string) (string, error) {
	// http://cmusphinx.sourceforge.net/wiki/faq
	// -ar 16000 sets frequency to required 16khz
	// -ac 1 sets the number of audio channels to 1
	newPath := filePath + "." + fileExt
	os.Remove(newPath) // If it already exists, ffmpeg will throw an error
	cmd := exec.Command("ffmpeg", "-i", filePath, "-ar", "16000", "-ac", "1", newPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", errors.New(err.Error() + "\nCommand Output:" + string(out))
	}
	return newPath, nil
}

// DownloadFileFromURL locally downloads an audio file stored at url.
func DownloadFileFromURL(url string) (string, error) {
	// Taken from https://github.com/thbar/golang-playground/blob/master/download-files.go
	filePath := filePathFromURL(url)
	file, err := os.Create(filePath)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer file.Close()

	// Get file contents
	response, err := http.Get(url)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer response.Body.Close()

	// Write the body to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", errors.Trace(err)
	}

	return filePath, nil
}

func filePathFromURL(url string) string {
	tokens := strings.Split(url, "/")
	filePath := tokens[len(tokens)-1]
	// ensure the filePath is unique by appending timestamp
	filePath = filePath + strconv.Itoa(int(time.Now().UnixNano()))
	return filePath
}

// MakeIBMTaskFunction returns a task function for transcription using IBM transcription functions.
func MakeIBMTaskFunction(audioURL string, emailAddresses []string, searchWords []string) (task func(string) error, onFailure func(string, string)) {
	task = func(id string) error {
		filePath, err := DownloadFileFromURL(audioURL)
		if err != nil {
			return errors.Trace(err)
		}
		defer os.Remove(filePath)

		log.WithField("task", id).
			Debugf("Downloaded file at %s to %s", audioURL, filePath)

		flacPath, err := ConvertAudioIntoFormat(filePath, "flac")
		if err != nil {
			return errors.Trace(err)
		}
		defer os.Remove(flacPath)

		log.WithField("task", id).
			Debugf("Converted file %s to %s", filePath, flacPath)

		ibmResult, err := TranscribeWithIBM(flacPath, config.Config.IBMUsername, config.Config.IBMPassword)
		if err != nil {
			return errors.Trace(err)
		}
		transcript := GetTranscript(ibmResult)

		log.WithField("task", id).
			Info(transcript)

		// TODO: save data to MongoDB and file to Backblaze.

		if err := SendEmail(config.Config.EmailUsername, config.Config.EmailPassword, "smtp.gmail.com", 25, emailAddresses, fmt.Sprintf("IBM Transcription %s Complete", id), "The transcript is below. It can also be found in the database."+"\n\n"+transcript); err != nil {
			return errors.Trace(err)
		}

		log.WithField("task", id).
			Debugf("Sent email to %v", emailAddresses)

		return nil
	}

	onFailure = func(id string, errMessage string) {
		err := SendEmail(config.Config.EmailUsername, config.Config.EmailPassword, "smtp.gmail.com", 25, emailAddresses, fmt.Sprintf("IBM Transcription %s Failed", id), errMessage)
		if err != nil {
			log.WithField("task", id).
				Debugf("Could not send error email to %v because of the error %v", emailAddresses, err.Error())
			return
		}
		log.WithField("task", id).
			Debugf("Sent email to %v", emailAddresses)
	}
	return task, onFailure
}
