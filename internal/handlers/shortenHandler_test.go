package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenHandler(t *testing.T) {
	type want struct {
		code         int
		contentType  string
		responseBody string
	}
	tests := []struct {
		name        string
		httpMethod  string
		requestPath string
		sourceURL   string
		shortID     string
		want        want
	}{
		{
			name:        "positive test #1",
			sourceURL:   "http://test.url",
			requestPath: "/api/shorten",
			httpMethod:  http.MethodPost,
			shortID:     "fVdpTFBo",
			want: want{
				code:         201,
				contentType:  HeaderApplicationJSON,
				responseBody: fmt.Sprintf("{\"result\":\"%v/fVdpTFBo\"}\n", baseURL),
			},
		}, {
			name:        "request url invalid",
			httpMethod:  http.MethodPost,
			requestPath: "/api/shorten",
			sourceURL:   "google.com",
			shortID:     "sKbYvAgT",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Full URL verify error\n",
			},
		}, {
			name:        "request body is empty",
			httpMethod:  http.MethodPost,
			requestPath: "/api/shorten",
			sourceURL:   "",
			shortID:     "sKbYvAgT",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Required url\n",
			},
		}, {
			name:        "http method not corrected",
			httpMethod:  http.MethodGet,
			requestPath: "/api/shorten",
			sourceURL:   "",
			shortID:     "sVpHyErn",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Http method not POST\n",
			},
		}, {
			name:        "internal server error",
			sourceURL:   "http://test6.url",
			requestPath: "/api/shorten",
			httpMethod:  http.MethodPost,
			shortID:     "fbUhNtPv",
			want: want{
				code:         500,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "random fail\n",
			},
		},
	}

	mockService := new(MockShorterService)
	handler := NewShortenHandler(mockService, baseURL)
	mockService.On("CreateShortLink", "http://test6.url").Return("", errors.New("random fail"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("CreateShortLink", tt.sourceURL).Return(tt.shortID, nil)

			requestData := ShortenRequest{
				URL: tt.sourceURL,
			}
			rqBytes, err := json.Marshal(requestData)
			if err != nil {
				assert.Error(t, err)
			}

			r := bytes.NewReader(rqBytes)
			postRequest := httptest.NewRequest(tt.httpMethod, tt.requestPath, r)
			postRequest.Header.Add("Content-Type", HeaderApplicationJSON)
			w := httptest.NewRecorder()
			handler.Handle(w, postRequest)
			res := w.Result()
			require.Equal(t, tt.want.code, res.StatusCode, "Incorrect status code")

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "no error for read response")
			err = res.Body.Close()
			assert.NoError(t, err, "no error for close response body")
			resBodyStr := string(resBody)

			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"), "Incorrect header content-type")
			assert.Equal(t, tt.want.responseBody, resBodyStr, "Incorrect response body")
		})
	}
}
