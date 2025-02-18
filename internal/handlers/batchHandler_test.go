package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladSnap/shortener/internal/services"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchHandler_Handle(t *testing.T) {
	type want struct {
		code         int
		contentType  string
		responseBody string
		shortedURLs  []ShortenRowResponse
	}
	tests := []struct {
		name         string
		httpMethod   string
		requestPath  string
		shortIDs     []string
		originalURLs []ShortenRowRequest
		want         want
	}{
		{
			name:        "positive test #1",
			requestPath: "/api/shorten/batch",
			httpMethod:  http.MethodPost,
			shortIDs:    []string{"hbFgvtUO", "fVdpTFBo"},
			originalURLs: []ShortenRowRequest{
				{CorrelationID: "crid1", OriginalURL: "http://test1.url/"},
				{CorrelationID: "crid2", OriginalURL: "http://test2.url/"}},
			want: want{
				code:        201,
				contentType: "application/json",
				shortedURLs: []ShortenRowResponse{
					{CorrelationID: "crid1", ShortURL: baseURL + "/hbFgvtUO"},
					{CorrelationID: "crid2", ShortURL: baseURL + "/fVdpTFBo"},
				},
			},
		}, {
			name:        "request url invalid",
			httpMethod:  http.MethodPost,
			requestPath: "/api/shorten/batch",
			shortIDs:    []string{"fvJhnBtG"},
			originalURLs: []ShortenRowRequest{
				{CorrelationID: "crid3", OriginalURL: "google.com"}},
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "OriginalURL verify error\n",
			},
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockShorterService(ctrl)
	batchHandler := NewBatchHandler(mockService, baseURL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var shortedLinks []*services.ShortedLink
			for i, s := range tt.shortIDs {
				shortedLinks = append(shortedLinks, &services.ShortedLink{
					URL:          s,
					CorelationID: tt.originalURLs[i].CorrelationID})
			}
			links := convertToShortLink(tt.originalURLs)
			mockService.EXPECT().CreateShortLinkBatch(ctx, links).
				Return(shortedLinks, nil).
				AnyTimes()

			var bufReq bytes.Buffer
			err := json.NewEncoder(&bufReq).Encode(tt.originalURLs)
			assert.NoError(t, err, "no error for encode json request")
			postRequest := httptest.NewRequest(tt.httpMethod, tt.requestPath, &bufReq)
			postRequest.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			batchHandler.Handle(w, postRequest)
			res := w.Result()
			require.Equal(t, tt.want.code, res.StatusCode, "Incorrect status code")
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"), "Incorrect header content-type")

			if res.Header.Get("Content-Type") == "application/json" {
				var shortedUrls []ShortenRowResponse
				err := json.NewDecoder(res.Body).Decode(&shortedUrls)
				assert.NoError(t, err, "no error for decode json response")
				isRepRowEq := sliceEquals(tt.want.shortedURLs, shortedUrls)
				assert.Equal(t, isRepRowEq, true, "response models not equal to want.shortedURLs")
			} else {
				resBody, err := io.ReadAll(res.Body)
				assert.NoError(t, err, "no error for read response")
				assert.Equal(t, tt.want.responseBody, string(resBody), "Incorrect response short url")
			}

			err = res.Body.Close()
			assert.NoError(t, err, "no error for close response body")
		})
	}
}

func sliceEquals[T ShortenRowResponse](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func convertToShortLink(requestRows []ShortenRowRequest) []*services.OriginalLink {
	links := make([]*services.OriginalLink, 0, len(requestRows))
	for _, r := range requestRows {
		lin := &services.OriginalLink{
			CorelationID: r.CorrelationID,
			URL:          r.OriginalURL,
		}
		links = append(links, lin)
	}

	return links
}
