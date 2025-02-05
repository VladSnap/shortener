package app

import (
	"fmt"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/data/repos"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/services"
)

func CreateServer(opts *config.Options, resMng *services.ResourceManager) (ShortenerServer, error) {
	var shortLinkRepo services.ShortLinkRepo

	if opts.FileStoragePath != "" {
		fileRepo, err := repos.NewFileShortLinkRepo(opts.FileStoragePath)
		if err != nil {
			return nil, fmt.Errorf("failed create FileShortLinkRepo: %w", err)
		}
		resMng.Register(fileRepo.Close)
		shortLinkRepo = fileRepo
	} else {
		shortLinkRepo = repos.NewShortLinkRepo()
	}

	shorterService := services.NewNaiveShorterService(shortLinkRepo)
	postHandler := handlers.NewPostHandler(shorterService, opts.BaseURL)
	getHandler := handlers.NewGetHandler(shorterService)
	shortenHandler := handlers.NewShortenHandler(shorterService, opts.BaseURL)
	pingHandler := handlers.NewGetPingHandler(opts)

	server := NewChiShortenerServer(opts, postHandler, getHandler, shortenHandler, pingHandler)
	return server, nil
}
