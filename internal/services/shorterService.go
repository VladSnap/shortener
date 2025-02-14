package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/helpers"
	"github.com/google/uuid"
)

type ShortLinkRepo interface {
	Add(ctx context.Context, link *data.ShortLinkData) (*data.ShortLinkData, error)
	AddBatch(ctx context.Context, links []*data.ShortLinkData) ([]*data.ShortLinkData, error)
	Get(ctx context.Context, shortID string) (*data.ShortLinkData, error)
}

// Генерирует мок для ShortLinkRepo
//go:generate mockgen -destination=mock_services_test.go -package services github.com/VladSnap/shortener/internal/services ShortLinkRepo

type NaiveShorterService struct {
	shortLinkRepo ShortLinkRepo
}

func NewNaiveShorterService(repo ShortLinkRepo) *NaiveShorterService {
	service := new(NaiveShorterService)
	service.shortLinkRepo = repo
	return service
}

const shortIDLength = 8

func (service *NaiveShorterService) CreateShortLink(ctx context.Context, originalURL string) (*ShortedLink, error) {
	id, shortID, err := createNewIds()
	if err != nil {
		return nil, fmt.Errorf("failed create ids: %w", err)
	}
	newLink := data.NewShortLinkData(id.String(), shortID, originalURL)
	createdLink, err := service.shortLinkRepo.Add(ctx, newLink)
	if err != nil {
		var duplErr *data.DuplicateShortLinkError
		if errors.As(err, &duplErr) {
			res := NewShortedLink("", "", "", duplErr.ShortURL, true)
			return res, nil
		}
		return nil, fmt.Errorf("failed create short link object: %w", err)
	}
	// Если короткие ссылки разные, значит был найден дубль и возвращено его значение.
	isDuplicate := shortID != createdLink.ShortURL
	res := NewShortedLink(createdLink.UUID, "", createdLink.OriginalURL, createdLink.ShortURL, isDuplicate)
	return res, nil
}

func (service *NaiveShorterService) GetURL(ctx context.Context, shortID string) (string, error) {
	link, err := service.shortLinkRepo.Get(ctx, shortID)
	if err != nil {
		return "", fmt.Errorf("failed get url from repo: %w", err)
	} else if link != nil {
		return link.OriginalURL, nil
	}
	return "", nil
}

func (service *NaiveShorterService) CreateShortLinkBatch(ctx context.Context, originalLinks []*OriginalLink) (
	[]*ShortedLink, error) {
	dataModels := make([]*data.ShortLinkData, 0, len(originalLinks))
	createdModels := make([]*ShortedLink, 0, len(originalLinks))

	if len(originalLinks) == 0 {
		return createdModels, nil
	}

	for _, ol := range originalLinks {
		id, shortID, err := createNewIds()
		if err != nil {
			return nil, fmt.Errorf("failed create ids: %w", err)
		}
		dm := data.NewShortLinkData(id.String(), shortID, ol.URL)
		dataModels = append(dataModels, dm)
		cm := NewShortedLink(id.String(), ol.CorelationID, ol.URL, shortID, false)
		createdModels = append(createdModels, cm)
	}

	_, err := service.shortLinkRepo.AddBatch(ctx, dataModels)
	if err != nil {
		return nil, fmt.Errorf("failed add batch in repo: %w", err)
	}
	// todo: Тут по хорошему надо обновить ShortURL в моделях, если в репозитории будет логика проверки дублей

	return createdModels, nil
}

func createNewIds() (id uuid.UUID, shortID string, err error) {
	id, err = uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("failed create random: %w", err)
		return
	}
	shortID, err = helpers.RandStringRunes(shortIDLength)
	if err != nil {
		err = fmt.Errorf("failed create short url: %w", err)
		return
	}
	return
}
