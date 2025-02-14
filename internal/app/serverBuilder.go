package app

import (
	"fmt"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/data/repos"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/services"
)

func CreateServer(opts *config.Options, resMng *services.ResourceManager) (ShortenerServer, error) {
	var shortLinkRepo services.ShortLinkRepo

	switch {
	case opts.DataBaseConnString != "":
		database, err := data.NewDatabaseShortener(opts.DataBaseConnString)
		if err != nil {
			return nil, fmt.Errorf("failed create DatabaseShortener: %w", err)
		}
		resMng.Register(database.Close)
		err = database.InitDatabase()
		if err != nil {
			return nil, fmt.Errorf("failed init Database: %w", err)
		}
		shortLinkRepo = repos.NewDatabaseShortLinkRepo(database)
	case opts.FileStoragePath != "":
		fileRepo, err := repos.NewFileShortLinkRepo(opts.FileStoragePath)
		if err != nil {
			return nil, fmt.Errorf("failed create FileShortLinkRepo: %w", err)
		}
		resMng.Register(fileRepo.Close)
		shortLinkRepo = fileRepo
	default:
		shortLinkRepo = repos.NewShortLinkRepo()
	}

	shorterService := services.NewNaiveShorterService(shortLinkRepo)
	postHandler := handlers.NewPostHandler(shorterService, opts.BaseURL)
	getHandler := handlers.NewGetHandler(shorterService)
	shortenHandler := handlers.NewShortenHandler(shorterService, opts.BaseURL)
	pingHandler := handlers.NewGetPingHandler(opts)
	batchHandler := handlers.NewBatchHandler(shorterService, opts.BaseURL)

	server := NewChiShortenerServer(opts, postHandler, getHandler, shortenHandler, pingHandler, batchHandler)
	return server, nil
}
