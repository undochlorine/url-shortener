package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Server    Server    `yaml:"Server" validate:"required" message:"required:{field} is required"`
		Shortener Shortener `yaml:"Shortener" validate:"required" message:"required:{field} is required"`
	}

	Server struct {
		Host     string `yaml:"Host" validate:"required" message:"required:{field} is required"`
		Port     string `yaml:"Port" validate:"required" message:"required:{field} is required"`
		LogLevel string `yaml:"LogLevel" validate:"required" message:"required:{field} is required (Info, Debug, Warn, Error)"`
	}

	Shortener struct {
		ConnString string `yaml:"ConnectionString" validate:"required" message:"required:{field} is required"`
	}
)

func New(path string) (*Config, error) {
	var config *Config

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c Config) valid() error {
	var err error
	if err = c.Server.valid(); err != nil {
		return err
	}
	if err = c.Shortener.valid(); err != nil {
		return err
	}
	return nil
}

func (s Server) valid() error {
	if s.Host == "" {
		return fmt.Errorf("[server] empty host name")
	}
	if s.Port == "" {
		return fmt.Errorf("[server] empty port")
	}
	if s.LogLevel == "" {
		return fmt.Errorf("[server] empty log level")
	}
	return nil
}

func (u Shortener) valid() error {
	if u.ConnString == "" {
		return fmt.Errorf("[shortener] empty connection string")
	}
	return nil
}
