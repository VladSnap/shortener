package main

import (
	"fmt"
	"os"

	"github.com/VladSnap/shortener/internal/app"
	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/config/validation"
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/handlers"
)

func main() {
	fmt.Println("Run shorneter server. Args:", os.Args)
	confValidator := &validation.OptionsValidator{}
	opts, err := config.LoadConfig(confValidator)
	if err != nil {
		panic(err)
	}

	server := createServer(opts)
	err = server.RunServer()

	if err != nil {
		panic(err)
	}
}

func createServer(opts *config.Options) app.ShortenerServer {
	shortLinkRepo := data.NewShortLinkRepo()
	getHandler := handlers.NewGetHandler(shortLinkRepo)
	postHandler := handlers.NewPostHandler(shortLinkRepo, opts.BaseURL)
	server := app.NewChiShortenerServer(opts, postHandler, getHandler)
	return server
}
