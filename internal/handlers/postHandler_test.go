package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockShortLinkRepo struct {
	mock.Mock
}

const baseUrl string = "http://localhost:8080/"

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
		sourceUrl   string
		shortId     string
		want        want
	}{
		{
			name:        "positive test #1",
			sourceUrl:   "http://test.url",
			requestPath: "/",
			httpMethod:  http.MethodPost,
			shortId:     "fVdpTFBo",
			want: want{
				code:         201,
				contentType:  "text/plain",
				responseBody: baseUrl + "fVdpTFBo",
			},
		}, {
			name:        "request body is empty",
			httpMethod:  http.MethodPost,
			requestPath: "/",
			sourceUrl:   "",
			shortId:     "sKbYvAgT",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
			},
		}, {
			name:        "http method not corrected",
			httpMethod:  http.MethodGet,
			requestPath: "/",
			sourceUrl:   "",
			shortId:     "sVpHyErn",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
			},
		}, {
			name:        "request path not correct #1",
			sourceUrl:   "http://test3.url",
			requestPath: "//",
			httpMethod:  http.MethodPost,
			shortId:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
			},
		}, {
			name:        "request path not correct #2",
			sourceUrl:   "http://test4.url",
			requestPath: "/foo",
			httpMethod:  http.MethodPost,
			shortId:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
			},
		}, {
			name:        "request path not correct #3",
			sourceUrl:   "http://test5.url",
			requestPath: "/foo/bar",
			httpMethod:  http.MethodPost,
			shortId:     "rDlUpOnb",
			want: want{
				code:         400,
				contentType:  "text/plain; charset=utf-8",
				responseBody: "Bad Request\n",
			},
		},
	}

	shortLinkRepo := new(MockShortLinkRepo)
	postHandler := NewPostHandler(shortLinkRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortLinkRepo.On("CreateShortLink", tt.sourceUrl).Return(tt.shortId)

			r := strings.NewReader(tt.sourceUrl)
			postRequest := httptest.NewRequest(tt.httpMethod, tt.requestPath, r)
			postRequest.Header.Add("Content-Type", "text/plain; charset=utf-8")
			w := httptest.NewRecorder()
			postHandler.Handle(w, postRequest)
			res := w.Result()
			require.Equal(t, tt.want.code, res.StatusCode, "Incorrect status code")
			resBody, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			shortUrl := string(resBody)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"), "Incorrect header content-type")
			assert.Equal(t, tt.want.responseBody, shortUrl, "Incorrect response short url")
		})
	}
}

func (repo *MockShortLinkRepo) CreateShortLink(url string) string {
	args := repo.Called(url)
	return args.String(0)
}

func (repo *MockShortLinkRepo) GetURL(key string) string {
	args := repo.Called(key)
	return args.String(0)
}
