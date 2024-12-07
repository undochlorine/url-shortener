package config

import (
	"errors"
	"fmt"
	"github.com/gookit/validate"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type (
	Config struct {
		LogLevel      string        `yaml:"LogLevel" validate:"required" message:"required:{field} is required (Info, Debug, Warn, Error)"`
		Host          string        `yaml:"Host" validate:"required" message:"required:{field} is required"`
		Database      string        `yaml:"Database" validate:"required" message:"required:{field} is required"`
		DBMaxIdleTime time.Duration `yaml:"DBMaxIdleTime"  validate:"required" message:"required:{field} is required"`
		AppMode       string        `yaml:"AppMode" validate:"required" message:"required:{field} is required (dev, stage, production)"`
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

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	if err = config.valid(); err != nil {
		return nil, err
	}

	return config, nil
}

const (
	APP_MODE_DEV        = "dev"
	APP_MODE_STAGE      = "stage"
	APP_MODE_PRODUCTION = "production"
)

func (s *Config) valid() error {
	if s.AppMode != APP_MODE_DEV && s.AppMode != APP_MODE_STAGE && s.AppMode != APP_MODE_PRODUCTION {
		return errors.New("invalid app mode")
	}

	v := validate.Struct(s)
	if v.Validate() {
		return nil
	}

	fmt.Println(v.Errors)
	return v.Errors.OneError()
}
