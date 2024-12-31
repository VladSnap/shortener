package config

import (
	flag "github.com/spf13/pflag"
)

type Options struct {
	ListenAddress string // server listen address
	BaseURL       string // base url for short url
}

func ParseFlags() *Options {
	opts := new(Options)

	flag.StringVar(&opts.ListenAddress, "a", ":8080", "server listen address")
	flag.StringVar(&opts.BaseURL, "b", "http://localhost:8080/", "base url for short url")

	flag.Parse()
	return opts
}
