package main

import (
	"fmt"
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

var resourceManager *services.ResourceManager

func main() {
	defer func() {
		err := log.Zap.Sync()
		panic(fmt.Errorf("failed zap logger sync: %w", err))
	}()
	resourceManager = services.NewResourceManager()
	defer func() {
		err := resourceManager.Cleanup()

		if err != nil {
			panic(fmt.Errorf("failed resourceManager clean: %w", err))
		}
	}()

	log.Zap.Info("run shorneter server", zap.Strings("Args", os.Args))

	confValidator := &validation.OptionsValidator{}
	opts, err := config.LoadConfig(confValidator)
	if err != nil {
		panic(err)
	}

	server, err := createServer(opts)
	if err != nil {
		panic(err)
	}

	err = server.RunServer()

	if err != nil {
		panic(err)
	}
}

func createServer(opts *config.Options) (app.ShortenerServer, error) {
	var shortLinkRepo services.ShortLinkRepo

	if opts.FileStoragePath != "" {
		fileRepo, err := data.NewFileShortLinkRepo(opts.FileStoragePath)
		if err != nil {
			return nil, fmt.Errorf("failed create FileShortLinkRepo: %w", err)
		}
		resourceManager.Register(fileRepo.Close)
		shortLinkRepo = fileRepo
	} else {
		shortLinkRepo = data.NewShortLinkRepo()
	}

	shorterService := services.NewNaiveShorterService(shortLinkRepo)
	postHandler := handlers.NewPostHandler(shorterService, opts.BaseURL)
	getHandler := handlers.NewGetHandler(shorterService)
	shortenHandler := handlers.NewShortenHandler(shorterService, opts.BaseURL)

	server := app.NewChiShortenerServer(opts, postHandler, getHandler, shortenHandler)
	return server, nil
}
