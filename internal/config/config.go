package config

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

type Options struct {
	ListenAddress string `env:"SERVER_ADDRESS"` // server listen address
	BaseURL       string `env:"BASE_URL"`       // base url for short url
}

func LoadConfig() *Options {
	opts := ParseFlags()
	ParseEnvConfig(opts)

	fmt.Printf("Config loaded: %+v\n", opts)
	return opts
}

func ParseFlags() *Options {
	opts := new(Options)

	flag.StringVar(&opts.ListenAddress, "a", ":8080", "server listen address")
	flag.StringVar(&opts.BaseURL, "b", "http://localhost:8080", "base url for short url")

	flag.Parse()

	runesBaseURL := []rune(opts.BaseURL)

	if string(runesBaseURL[len(runesBaseURL)-1:]) == "/" {
		panic("Incorrect -b argument. Don't put a slash at the end")
	}

	return opts
}

func ParseEnvConfig(opts *Options) {
	err := env.Parse(opts)

	if err != nil {
		log.Fatal(err)
	}
}
