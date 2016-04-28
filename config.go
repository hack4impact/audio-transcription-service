package main

import (
	"github.com/BurntSushi/toml"
)

//parseConfigFile parses the specified file into a struct
func parseConfigFile(filename string) (*AppConfig, error) {
	var config AppConfig
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Config is an example
type AppConfig struct {
	EmailUsername         string
	EmailPassword         string
	TranscriptionServices []string
}
