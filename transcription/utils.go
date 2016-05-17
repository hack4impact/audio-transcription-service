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

// SplitFlacFile ensures that the input audio files to IBM are less than 100mb, with 5 seconds of redundancy between files.
func SplitFlacFile(fn string) error {
	// http://stackoverflow.com/questions/36632511/split-audio-file-into-several-files-each-below-a-size-threshold
	// The Stack Overflow answer ultimately calculated the length of each audio chunk in seconds.
	// chunk_length_in_sec = math.ceil((duration_in_sec * file_split_size ) / wav_file_size)
	// Invariant: If ConvertAudioIntoWavFormat is called on fn, a 95MB chunk of resulting Wav file is always 2968 seconds.
	// In the above equation, there is one constant: file_split_size = 95000000 bytes.
	// duration_in_sec is used to calculate wav_file_size, so it is canceled out in the ratio.
	// wav_file_size = (sample_rate * bit_rate * channel_count * duration_in_sec) / 8
	// sample_rate = 44100, bit_rate = 16, channels_count = 1 (stereo: 2, but Sphinx prefers 1)
	// As a chunk of the Wav file is extracted using FFMPEG, it is converted back into Flac format.
	wavFileName, err := ConvertAudioIntoFormat(fn, "wav")
	if err != nil {
		return err
	}

	numChunks, err := getNumChunks(fn)
	if err != nil {
		return err
	}

	chunkLengthInSeconds := 2968
	for i := 0; i < numChunks; i++ {
		// 5 seconds of redundancy for each chunk after the first
		startingSecond := i*chunkLengthInSeconds - (i-1)*5
		newFileName := wavFileName + strconv.Itoa(i)
		if err := extractAudioSegment(newFileName, startingSecond, chunkLengthInSeconds); err != nil {
			return err
		}
		if _, err := ConvertAudioIntoFormat(newFileName, "flac"); err != nil {
			return err
		}
	}
	return nil
}

// getNumChunks gets file size in MB, divides by 95 MB, and add 1 more chunk in case
func getNumChunks(fn string) (int, error) {
	file, err := os.Open(fn)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return -1, err
	}

	wavFileSize := int(stat.Size())
	fileSplitSize := 95000000
	// The redundant seconds (5 seconds for every ~50 mintues) won't add own chunk
	// In case the remainder is almost the file size, add one more chunk
	numChunks := wavFileSize/fileSplitSize + 1
	return numChunks, nil
}

// extractAudioSegment uses FFMPEG to write a new audio file starting at a given time of a given length
func extractAudioSegment(fn string, ss int, t int) error {
	// -ss: starting second, -t: time in seconds
	cmd := exec.Command("ffmpeg", "-i", fn, "-ss", strconv.Itoa(ss), "-t", strconv.Itoa(t), fn)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
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
