package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockShorterService struct {
	mock.Mock
}

const baseURL string = "http://localhost:8080"

func TestPostHandler(t *testing.T) {
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
			requestPath: "/",
			httpMethod:  http.MethodPost,
			shortID:     "fVdpTFBo",
			want: want{
				code:         201,
				contentType:  "text/plain",
				responseBody: baseURL + "/fVdpTFBo",
			},
		}, {
			name:        "request url invalid",
			httpMethod:  http.MethodPost,
			requestPath: "/",
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
			requestPath: "/",
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
			requestPath: "/",
			sourceURL:   "",
			shortID:     "sVpHyErn",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Http method not POST\n",
			},
		}, {
			name:        "request path not correct #1",
			sourceURL:   "http://test3.url",
			requestPath: "//",
			httpMethod:  http.MethodPost,
			shortID:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Incorrect request path\n",
			},
		}, {
			name:        "request path not correct #2",
			sourceURL:   "http://test4.url",
			requestPath: "/foo",
			httpMethod:  http.MethodPost,
			shortID:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Incorrect request path\n",
			},
		}, {
			name:        "request path not correct #3",
			sourceURL:   "http://test5.url",
			requestPath: "/foo/bar",
			httpMethod:  http.MethodPost,
			shortID:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Incorrect request path\n",
			},
		}, {
			name:        "internal server error",
			sourceURL:   "http://test6.url",
			requestPath: "/",
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
	postHandler := NewPostHandler(mockService, baseURL)
	mockService.On("CreateShortLink", "http://test6.url").Return("", errors.New("random fail"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("CreateShortLink", tt.sourceURL).Return(tt.shortID, nil)

			r := strings.NewReader(tt.sourceURL)
			postRequest := httptest.NewRequest(tt.httpMethod, tt.requestPath, r)
			postRequest.Header.Add("Content-Type", "text/plain; charset=utf-8")
			w := httptest.NewRecorder()
			postHandler.Handle(w, postRequest)
			res := w.Result()
			require.Equal(t, tt.want.code, res.StatusCode, "Incorrect status code")
			resBody, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			shortURL := string(resBody)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"), "Incorrect header content-type")
			assert.Equal(t, tt.want.responseBody, shortURL, "Incorrect response short url")
		})
	}
}

func (repo *MockShorterService) CreateShortLink(url string) (string, error) {
	args := repo.Called(url)
	return args.String(0), args.Error(1)
}

func (repo *MockShorterService) GetURL(key string) string {
	args := repo.Called(key)
	return args.String(0)
}
