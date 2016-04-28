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

// ConvertAudioIntoWavFormat converts encoded audio into the required format.
func ConvertAudioIntoWavFormat(fn string) error {
	// http://cmusphinx.sourceforge.net/wiki/faq
	// -ar 16000 sets frequency to required 16khz
	// -ac 1 sets the number of audio channels to 1
	cmd := exec.Command("ffmpeg", "-i", fn, "-ar", "16000", "-ac", "1", fn+".wav")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// ConvertAudioIntoFlacFormat converts files into .flac format.
func ConvertAudioIntoFlacFormat(fn string) error {
	// -ar 16000 sets frequency to required 16khz
	// -ac 1 sets the number of audio channels to 1
	cmd := exec.Command("ffmpeg", "-i", fn, "-ar", "16000", "-ac", "1", fn+".flac")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// DownloadFileFromURL locally downloads an audio file stored at url.
func DownloadFileFromURL(url string) error {
	// Taken from https://github.com/thbar/golang-playground/blob/master/download-files.go
	output, err := os.Create(fileNameFromURL(url))
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

func fileNameFromURL(url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	return fileName
}

// SplitFlacFile ensures that the input audio files to IBM are less than 100mb.
func SplitFlacFile(fn string) error {
	// http://stackoverflow.com/questions/36632511/split-audio-file-into-several-files-each-below-a-size-threshold
	// The answer ultimately calculates the length of each audio chunk in seconds.
	// chunk_length_in_sec = math.ceil((duration_in_sec * file_split_size ) / wav_file_size)
	// If ConvertAudioIntoWavFormat is called on fn, a 95MB chunk is always 2968 seconds.
	// In the above equation, there is one constant: file_split_size = 95000000 bytes.
	// duration_in_sec is used to calculate wav_file_size, so it is canceled out in the ratio.
	// wav_file_size = (sample_rate * bit_rate * channel_count * duration_in_sec) / 8
	// sample_rate = 44100, bit_rate = 16, channels_count = 1 (stereo: 2, Sphinx: 1)
	// Afterwards, every Wav file is converted back into Flac format.
	err := ConvertAudioIntoWavFormat(fn)
	if err != nil {
		return err
	}
	wavFileName := strings.Split(fn, ".")[0] + ".wav"

	numChunks, err := getNumChunks(fn)
	if err != nil {
		return err
	}

	chunkTime := 2968
	ss := 0
	for i := 0; i < numChunks; i++ {
		if err := extractAudioSegment(wavFileName+strconv.Itoa(i), ss, chunkTime); err != nil {
			return err
		}
		if err := ConvertAudioIntoFlacFormat(wavFileName + strconv.Itoa(i)); err != nil {
			return err
		}
		ss -= 5
	}
	return nil
}

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
	wavFileSize := stat.Size()
	numChunks := 0
	fileSplitSize := 95000000
	fiveSecondsInBytes := 16000 // 5s / 2968s * 9500000 bytes
	bytesUsed := 0
	for bytesUsed < wavFileSize {
		numChunks++
		bytesUsed += int64(fileSplitSize) - int64(fiveSecondsInBytes)
	}
}

func extractAudioSegment(fn string, ss int, t int) error {
	// -ss: starting second, -t: time in seconds
	cmd := exec.Command("ffmpeg", "-i", fn, "-ss", strconv.Itoa(ss), "-t", strconv.Itoa(t), fn)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// MakeTaskFunction returns a task function for transcription using transcription functions.
func MakeTaskFunction(audioURL string, emailAddresses []string, searchWords []string) func() error {
	return func() error {
		fileName := fileNameFromURL(audioURL)
		if err := DownloadFileFromURL(audioURL); err != nil {
			return err
		}
		if err := ConvertAudioIntoWavFormat(fileName); err != nil {
			return err
		}
		return nil
	}
}
