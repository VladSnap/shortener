package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockShortLinkRepo struct {
	mock.Mock
}

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

	mockRepo := new(MockShortLinkRepo)
	service := NewNaiveShorterService(mockRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("CreateShortLink", mock.AnythingOfType("string"), tt.sourceURL).Return(nil)
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

	mockRepo := new(MockShortLinkRepo)
	service := NewNaiveShorterService(mockRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetURL", tt.shortID).Return(tt.want.fullURL)
			result := service.GetURL(tt.shortID)

			assert.Equal(t, tt.want.fullURL, result)
		})
	}
}

func (repo *MockShortLinkRepo) CreateShortLink(shortID string, fullURL string) error {
	args := repo.Called(shortID, fullURL)
	return args.Error(0)
}

func (repo *MockShortLinkRepo) GetURL(key string) string {
	args := repo.Called(key)
	return args.String(0)
}
