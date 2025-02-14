package config

import (
	"flag"
	"fmt"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/caarlos0/env/v6"
)

type Options struct {
	ListenAddress      string `env:"SERVER_ADDRESS"`    // server listen address
	BaseURL            string `env:"BASE_URL"`          // base url for short url
	FileStoragePath    string `env:"FILE_STORAGE_PATH"` // file path to storage all shorten url
	DataBaseConnString string `env:"DATABASE_DSN"`      // database connection string
}

type ConfigValidater interface {
	Validate(opts *Options) error
}

func LoadConfig(validater ConfigValidater) (*Options, error) {
	opts, err := ParseFlags(validater)
	if err != nil {
		return nil, err
	}
	err = ParseEnvConfig(opts)
	if err != nil {
		return nil, err
	}

	log.Zap.Infof("Config loaded: %+v\n", opts)
	return opts, nil
}

func ParseFlags(validater ConfigValidater) (*Options, error) {
	opts := new(Options)

	flag.StringVar(&opts.ListenAddress, "a", ":8080", "server listen address")
	flag.StringVar(&opts.BaseURL, "b", "http://localhost:8080", "base url for short url")
	flag.StringVar(&opts.FileStoragePath, "f", "", "file path to storage all shorten url")
	flag.StringVar(&opts.DataBaseConnString, "d", "", "database connection string")

	flag.Parse()

	err := validater.Validate(opts)

	if err != nil {
		return nil, fmt.Errorf("config validating failed: %w", err)
	}

	return opts, nil
}

func ParseEnvConfig(opts *Options) error {
	err := env.Parse(opts)

	if err != nil {
		return fmt.Errorf("failed env parsing: %w", err)
	}

	return nil
}
