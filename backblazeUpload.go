package main

import (
	"os"
	"path/filepath"

	"gopkg.in/kothar/go-backblaze.v0"
)

func uploadFileToBackblaze(filename string, accountID string, applicationKey string, bucketName string) (string, error) {
	b2, err := backblaze.NewB2(backblaze.Credentials{
		AccountID:      accountID,
		ApplicationKey: applicationKey,
	})
	if err != nil {
		return "", err
	}

	bucket, err := b2.Bucket(bucketName)
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
