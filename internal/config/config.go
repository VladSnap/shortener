package config

import (
	"flag"
)

type Options struct {
	ListenAddress string // server listen address
	BaseURL       string // base url for short url
}

func ParseFlags() *Options {
	opts := new(Options)

	flag.StringVar(&opts.ListenAddress, "a", ":8080", "server listen address")
	flag.StringVar(&opts.BaseURL, "b", "http://localhost:8080", "base url for short url")

	flag.Parse()

	runesBaseUrl := []rune(opts.BaseURL)

	if string(runesBaseUrl[len(runesBaseUrl)-1:]) == "/" {
		panic("Incorrect -b argument. Don't put a slash at the end")
	}

	return opts
}
