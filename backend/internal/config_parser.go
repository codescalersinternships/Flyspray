package internal

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/validator.v2"
)

// Configuration is app configuration data
type Configuration struct {
	Server     Server     `json:"server"`
	MailSender MailSender `json:"mail_sender"`
	DB         DB         `json:"db"`
	JWT        JWT        `json:"jwt"`
	Version    Version    `json:"version"`
}

// Server is server configuration data
type Server struct {
	Host string `json:"host" validate:"nonzero"`
	Port int    `json:"port" validate:"nonzero"`
}

// MailSender is mail sender configuration data
type MailSender struct {
	Email       string `json:"email" validate:"nonzero"`
	SendGridKey string `json:"sendgrid_key" validate:"nonzero"`
	Timeout     int    `json:"timeout" validate:"min=30"`
}

// DB is database configuration data
type DB struct {
	File string `json:"file" validate:"nonzero"`
}

// JWT is jwt configuration data
type JWT struct {
	Secret  string `json:"secret" validate:"nonzero"`
	Timeout int    `json:"timeout" validate:"min=5"`
}

// Version represents the version of the app
type Version struct {
	AppVersion string `json:"app_version"`
}

// ReadConfigFile read configuration from json file
func ReadConfigFile(path string) (Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return Configuration{}, fmt.Errorf("failed to open config file: %w", err)
	}

	config := Configuration{}
	dec := json.NewDecoder(file)

	if err := dec.Decode(&config); err != nil {
		return Configuration{}, fmt.Errorf("failed to load config file: %w", err)
	}

	return config, validator.Validate(config)
}
