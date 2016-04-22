package main

import (
	"os"
	"path/filepath"

	"gopkg.in/kothar/go-backblaze.v0"
)

func uploadFileToBackblaze(filename string) (string, error) {
	b2, err := backblaze.NewB2(backblaze.Credentials{
		AccountID:      "23547fcec776",
		ApplicationKey: "0016ab4da23ef8548aa6d19c77e0eada59ae55764e",
	})
	if err != nil {
		return "", err
	}

	bucket, err := b2.Bucket("Hack4Impact")
	if err != nil {
		return "", err
	}

	path := filename
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}

	name := filepath.Base(path)
	metadata := make(map[string]string)

	_, err = bucket.UploadFile(name, metadata, reader)
	if err != nil {
		return "", err
	}

	return bucket.FileURL(name)
}
