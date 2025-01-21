package services

import (
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/helpers"
)

type ShorterService interface {
	CreateShortLink(fullURL string) (string, error)
	GetURL(shortID string) string
}

type NaiveShorterService struct {
	shortLinkRepo data.ShortLinkRepo
}

func NewNaiveShorterService(repo data.ShortLinkRepo) *NaiveShorterService {
	service := new(NaiveShorterService)
	service.shortLinkRepo = repo
	return service
}

const shortIDLength = 8

func (service *NaiveShorterService) CreateShortLink(fullURL string) (string, error) {
	shortID, err := helpers.RandStringRunes(shortIDLength)
	if err != nil {
		return "", err
	}

	err = service.shortLinkRepo.CreateShortLink(shortID, fullURL)
	if err != nil {
		return "", err
	}

	return shortID, nil
}

func (service *NaiveShorterService) GetURL(shortID string) string {
	return service.shortLinkRepo.GetURL(shortID)
}
