// Package transcription implements functions for the manipulation and
// transcription of audio files.
package transcription

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
)

//Transcription contains the transcription text and metadata of a transcription
//job
type Transcription struct {
	TextTranscription string
	Metadata          string
}

// SphinxTranscription transcribes a given file using Sphinx.
// File name should not include the type extension.
func SphinxTranscription(fileName string) (Transcription, error) {
	var result Transcription
	os.Chdir("./Sphinx")
	p, err := os.Getwd()
	if err != nil {
		return result, err
	}
	cmd := exec.Command("bash", p+"/gradlew", "run", "-Pmyargs=files/"+fileName)
	if err := cmd.Run(); err != nil {
		return result, err
	}
	outputFile := p + "files/" + fileName + "-json.txt"
	result, err = outputToStruct(outputFile)
	if err != nil {
		return result, err
	}
	return result, nil
}

// outputToStruct takes a text file with json data and reads its input
// into a Go struct
func outputToStruct(fileName string) (Transcription, error) {
	var jsonData Transcription
	file, err := os.Open(fileName)
	if err != nil {
		return jsonData, err
	}
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)
	err = dec.Decode(&jsonData)
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

//Need to add a corresponding MakeSphinxTaskFunction
