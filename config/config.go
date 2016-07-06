package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

// Config is the application-wide config
var Config AppConfig

func init() {
	err := parseConfigFile(&Config, "config.toml")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%+v", Config)
}

// parseConfigFile parses the specified file into a given struct
func parseConfigFile(config *AppConfig, filename string) error {
	if _, err := toml.DecodeFile(filename, config); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// AppConfig contains the app config variables.
type AppConfig struct {
	BackblazeAccountID      string
	BackblazeApplicationKey string
	BackblazeBucket         string
	Debug                   bool
	EmailUsername           string
	EmailPassword           string
	IBMUsername             string
	IBMPassword             string
	MongoURL                string
	Port                    int
	SecretKey               string
}
