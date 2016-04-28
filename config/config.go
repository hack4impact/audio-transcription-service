package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

// Config is the application-wide config
var Config AppConfig

func init() {
	err := parseConfigFile(&Config, "config.toml")
	if err != nil {
		log.Fatal(err)
	}
	// log.Info(Config)
	log.Infof("%+v", Config)
}

//parseConfigFile parses the specified file into a given struct
func parseConfigFile(config *AppConfig, filename string) error {
	if _, err := toml.DecodeFile(filename, config); err != nil {
		return err
	}
	return nil
}

// AppConfig contains the app config variables.
type AppConfig struct {
	EmailUsername string
	EmailPassword string
	IBMUsername   string
	IBMPassword   string
	Debug         bool
}
