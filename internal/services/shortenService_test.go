package services

import (
	"context"
	"testing"

	"github.com/VladSnap/shortener/internal/data"
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

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockShortLinkRepo(ctrl)
	service := NewNaiveShorterService(mockRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retLink := getNewShortLink("tttttttt", tt.sourceURL)
			mockRepo.EXPECT().Add(ctx, gomock.Any()).Return(retLink, nil)
			result, err := service.CreateShortLink(context.Background(), tt.sourceURL)

			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.URL)
			assert.Len(t, result.URL, shortIDLength)
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

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockShortLinkRepo(ctrl)
	service := NewNaiveShorterService(mockRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retLink := getNewShortLink(tt.shortID, tt.want.fullURL)
			mockRepo.EXPECT().Get(ctx, tt.shortID).Return(retLink, nil)
			result, err := service.GetURL(context.Background(), tt.shortID)
			assert.NoError(t, err, "no expect error get url")
			assert.Equal(t, tt.want.fullURL, result)
		})
	}
}

func getNewShortLink(shortID string, originalURL string) *data.ShortLinkData {
	id := uuid.MustParse("2093ad7c-6227-4d97-8f83-9e837ab6474b")
	return &data.ShortLinkData{UUID: id.String(), ShortURL: shortID, OriginalURL: originalURL}
}
