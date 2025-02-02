package services

import (
	"fmt"

	"github.com/VladSnap/shortener/internal/helpers"
)

type ShortLinkRepo interface {
	CreateShortLink(shortID string, fullURL string) error
	GetURL(shortID string) string
}

type ShorterService interface {
	CreateShortLink(fullURL string) (string, error)
	GetURL(shortID string) string
}

type NaiveShorterService struct {
	shortLinkRepo ShortLinkRepo
}

func NewNaiveShorterService(repo ShortLinkRepo) *NaiveShorterService {
	service := new(NaiveShorterService)
	service.shortLinkRepo = repo
	return service
}

const shortIDLength = 8

func (service *NaiveShorterService) CreateShortLink(fullURL string) (string, error) {
	shortID, err := helpers.RandStringRunes(shortIDLength)
	if err != nil {
		return "", fmt.Errorf("failed create short url: %w", err)
	}

	err = service.shortLinkRepo.CreateShortLink(shortID, fullURL)
	if err != nil {
		return "", fmt.Errorf("failed create short link object: %w", err)
	}

	return shortID, nil
}

func (service *NaiveShorterService) GetURL(shortID string) string {
	return service.shortLinkRepo.GetURL(shortID)
}
