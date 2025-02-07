package services

import (
	"testing"

	"github.com/VladSnap/shortener/internal/data/models"
	"github.com/VladSnap/shortener/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNaiveShortenService_CreateShortLink(t *testing.T) {
	tests := []struct {
		name      string
		sourceURL string
	}{
		{
			name:      "CreateShortLink positive test#1",
			sourceURL: "http://test.url",
		},
	}

	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// создаём объект-заглушку
	mockRepo := mocks.NewMockShortLinkRepo(ctrl)
	service := NewNaiveShorterService(mockRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retLink := getNewShortLink("tttttttt", tt.sourceURL)
			mockRepo.EXPECT().CreateShortLink(gomock.Any()).Return(retLink, nil)
			result, err := service.CreateShortLink(tt.sourceURL)

			assert.Nil(t, err)
			assert.NotEmpty(t, result)
			assert.Len(t, result, shortIDLength)
		})
	}
}

func TestNaiveShortenService_GetURL(t *testing.T) {
	type want struct {
		fullURL string
	}
	tests := []struct {
		name    string
		shortID string
		want    want
	}{
		{
			name:    "GetURL positive test#1",
			shortID: "avFjNyBR",
			want: want{
				fullURL: "http://test.url",
			},
		}, {
			name:    "GetURL not found shortID",
			shortID: "notFOUND",
			want: want{
				fullURL: "",
			},
		},
	}

	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// создаём объект-заглушку
	mockRepo := mocks.NewMockShortLinkRepo(ctrl)
	service := NewNaiveShorterService(mockRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retLink := getNewShortLink(tt.shortID, tt.want.fullURL)
			mockRepo.EXPECT().GetURL(tt.shortID).Return(retLink, nil)
			result, err := service.GetURL(tt.shortID)
			assert.NoError(t, err, "no expect error get url")
			assert.Equal(t, tt.want.fullURL, result)
		})
	}
}

func getNewShortLink(shortID string, originalURL string) *models.ShortLinkData {
	id := uuid.MustParse("2093ad7c-6227-4d97-8f83-9e837ab6474b")
	return &models.ShortLinkData{UUID: id.String(), ShortURL: shortID, OriginalURL: originalURL}
}
