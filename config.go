package main

import (
	"github.com/BurntSushi/toml"
)

//parseConfigFile parses the specified file into a struct
func parseConfigFile(filename string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Config is an example
type Config struct {
	EmailUsername  string
	EmailPassword  string
	AccountID      string
	ApplicationKey string
	BucketName     string
}
