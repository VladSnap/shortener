package main

import (
	"os"

	"github.com/VladSnap/shortener/internal/app"
	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/config/validation"
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
	"go.uber.org/zap"
)

func main() {
	defer log.Zap.Sync()

	log.Zap.Info("run shorneter server", zap.Strings("Args", os.Args))

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
	shorterService := services.NewNaiveShorterService(shortLinkRepo)
	getHandler := handlers.NewGetHandler(shorterService)
	postHandler := handlers.NewPostHandler(shorterService, opts.BaseURL)
	server := app.NewChiShortenerServer(opts, postHandler, getHandler)
	return server
}
