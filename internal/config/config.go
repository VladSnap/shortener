package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Options struct {
	ListenAddress string `env:"SERVER_ADDRESS"` // server listen address
	BaseURL       string `env:"BASE_URL"`       // base url for short url
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

	fmt.Printf("Config loaded: %+v\n", opts)
	return opts, nil
}

func ParseFlags(validater ConfigValidater) (*Options, error) {
	opts := new(Options)

	flag.StringVar(&opts.ListenAddress, "a", ":8080", "server listen address")
	flag.StringVar(&opts.BaseURL, "b", "http://localhost:8080", "base url for short url")

	flag.Parse()

	err := validater.Validate(opts)

	return opts, err
}

func ParseEnvConfig(opts *Options) error {
	err := env.Parse(opts)

	if err != nil {
		return err
	}

	return nil
}
